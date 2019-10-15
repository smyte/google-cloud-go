/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/smyte/google-cloud-go/internal/trace"
	"github.com/smyte/google-cloud-go/spanner/internal/backoff"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	edpb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const (
	retryInfoKey = "google.rpc.retryinfo-bin"
)

// errRetry returns an unavailable error under error namespace EsOther. It is a
// generic retryable error that is used to mask and recover unretryable errors
// in a retry loop.
func errRetry(err error) error {
	if se, ok := err.(*Error); ok {
		return &Error{codes.Unavailable, fmt.Sprintf("generic Cloud Spanner retryable error: { %v }", se.Error()), se.trailers}
	}
	return spannerErrorf(codes.Unavailable, "generic Cloud Spanner retryable error: { %v }", err.Error())
}

// isErrorClosing reports whether the error is generated by gRPC layer talking
// to a closed server.
func isErrorClosing(err error) bool {
	if err == nil {
		return false
	}
	if ErrCode(err) == codes.Internal && strings.Contains(ErrDesc(err), "transport is closing") {
		// Handle the case when connection is closed unexpectedly.
		// TODO: once gRPC is able to categorize this as retryable error, we
		//  should stop parsing the error message here.
		return true
	}
	return false
}

// isErrorRST reports whether the error is generated by gRPC client receiving a
// RST frame from server.
func isErrorRST(err error) bool {
	if err == nil {
		return false
	}
	if ErrCode(err) == codes.Internal && strings.Contains(ErrDesc(err), "stream terminated by RST_STREAM") {
		// TODO: once gRPC is able to categorize this error as "go away" or "retryable",
		// we should stop parsing the error message.
		return true
	}
	return false
}

// isErrorUnexpectedEOF returns true if error is generated by gRPC layer
// receiving io.EOF unexpectedly.
func isErrorUnexpectedEOF(err error) bool {
	if err == nil {
		return false
	}
	// Unexpected EOF is a transport layer issue that could be recovered by
	// retries. The most likely scenario is a flaky RecvMsg() call due to
	// network issues.
	//
	// For grpc version >= 1.14.0, the error code is Internal.
	// (https://github.com/grpc/grpc-go/releases/tag/v1.14.0)
	if ErrCode(err) == codes.Internal && strings.Contains(ErrDesc(err), "unexpected EOF") {
		return true
	}
	// For grpc version < 1.14.0, the error code in Unknown.
	if ErrCode(err) == codes.Unknown && strings.Contains(ErrDesc(err), "unexpected EOF") {
		return true
	}
	return false
}

// isErrorUnavailable returns true if the error is about server being
// unavailable.
func isErrorUnavailable(err error) bool {
	if err == nil {
		return false
	}
	if ErrCode(err) == codes.Unavailable {
		return true
	}
	return false
}

// isRetryable returns true if the Cloud Spanner error being checked is a
// retryable error.
func isRetryable(err error) bool {
	if isErrorClosing(err) {
		return true
	}
	if isErrorUnexpectedEOF(err) {
		return true
	}
	if isErrorRST(err) {
		return true
	}
	if isErrorUnavailable(err) {
		return true
	}
	return false
}

// errContextCanceled returns *spanner.Error for canceled context.
func errContextCanceled(ctx context.Context, lastErr error) error {
	if ctx.Err() == context.DeadlineExceeded {
		return spannerErrorf(codes.DeadlineExceeded, "%v, lastErr is <%v>", ctx.Err(), lastErr)
	}
	return spannerErrorf(codes.Canceled, "%v, lastErr is <%v>", ctx.Err(), lastErr)
}

// extractRetryDelay extracts retry backoff if present.
func extractRetryDelay(err error) (time.Duration, bool) {
	trailers := errTrailers(err)
	if trailers == nil {
		return 0, false
	}
	elem, ok := trailers[retryInfoKey]
	if !ok || len(elem) <= 0 {
		return 0, false
	}
	_, b, err := metadata.DecodeKeyValue(retryInfoKey, elem[0])
	if err != nil {
		return 0, false
	}
	var retryInfo edpb.RetryInfo
	if proto.Unmarshal([]byte(b), &retryInfo) != nil {
		return 0, false
	}
	delay, err := ptypes.Duration(retryInfo.RetryDelay)
	if err != nil {
		return 0, false
	}
	return delay, true
}

// runRetryable keeps attempting to run f until one of the following happens:
//     1) f returns nil error or an unretryable error;
//     2) context is cancelled or timeout.
//
// TODO: consider using https://github.com/googleapis/gax-go/v2 once it
// becomes available internally.
func runRetryable(ctx context.Context, f func(context.Context) error) error {
	return toSpannerError(runRetryableNoWrap(ctx, f))
}

// Like runRetryable, but doesn't wrap the returned error in a spanner.Error.
func runRetryableNoWrap(ctx context.Context, f func(context.Context) error) error {
	var funcErr error
	retryCount := 0
	for {
		select {
		case <-ctx.Done():
			// Do context check here so that even f() failed to do so (for
			// example, gRPC implementation bug), the loop can still have a
			// chance to exit as expected.
			return errContextCanceled(ctx, funcErr)
		default:
		}
		funcErr = f(ctx)
		if funcErr == nil {
			return nil
		}
		if isRetryable(funcErr) {
			// Error is retryable, do exponential backoff and continue.
			b, ok := extractRetryDelay(funcErr)
			if !ok {
				b = backoff.DefaultBackoff.Delay(retryCount)
			}
			trace.TracePrintf(ctx, nil, "Backing off for %s, then retrying", b)
			select {
			case <-ctx.Done():
				return errContextCanceled(ctx, funcErr)
			case <-time.After(b):
			}
			retryCount++
			continue
		}
		// Error isn't retryable / no error, return immediately.
		return funcErr
	}
}
