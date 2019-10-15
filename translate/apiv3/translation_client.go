// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gapic-generator. DO NOT EDIT.

package translate

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/smyte/google-cloud-go/longrunning"
	lroauto "github.com/smyte/google-cloud-go/longrunning/autogen"
	"github.com/golang/protobuf/proto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// TranslationCallOptions contains the retry settings for each method of TranslationClient.
type TranslationCallOptions struct {
	TranslateText         []gax.CallOption
	DetectLanguage        []gax.CallOption
	GetSupportedLanguages []gax.CallOption
	BatchTranslateText    []gax.CallOption
	CreateGlossary        []gax.CallOption
	ListGlossaries        []gax.CallOption
	GetGlossary           []gax.CallOption
	DeleteGlossary        []gax.CallOption
}

func defaultTranslationClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("translate.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultTranslationCallOptions() *TranslationCallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.3,
				})
			}),
		},
	}
	return &TranslationCallOptions{
		TranslateText:         retry[[2]string{"default", "non_idempotent"}],
		DetectLanguage:        retry[[2]string{"default", "non_idempotent"}],
		GetSupportedLanguages: retry[[2]string{"default", "idempotent"}],
		BatchTranslateText:    retry[[2]string{"default", "non_idempotent"}],
		CreateGlossary:        retry[[2]string{"default", "non_idempotent"}],
		ListGlossaries:        retry[[2]string{"default", "idempotent"}],
		GetGlossary:           retry[[2]string{"default", "idempotent"}],
		DeleteGlossary:        retry[[2]string{"default", "non_idempotent"}],
	}
}

// TranslationClient is a client for interacting with Cloud Translation API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type TranslationClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	translationClient translatepb.TranslationServiceClient

	// LROClient is used internally to handle longrunning operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient *lroauto.OperationsClient

	// The call options for this service.
	CallOptions *TranslationCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewTranslationClient creates a new translation service client.
//
// Provides natural language translation operations.
func NewTranslationClient(ctx context.Context, opts ...option.ClientOption) (*TranslationClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultTranslationClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &TranslationClient{
		conn:        conn,
		CallOptions: defaultTranslationCallOptions(),

		translationClient: translatepb.NewTranslationServiceClient(conn),
	}
	c.setGoogleClientInfo()

	c.LROClient, err = lroauto.NewOperationsClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		// This error "should not happen", since we are just reusing old connection
		// and never actually need to dial.
		// If this does happen, we could leak conn. However, we cannot close conn:
		// If the user invoked the function with option.WithGRPCConn,
		// we would close a connection that's still in use.
		// TODO(pongad): investigate error conditions.
		return nil, err
	}
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *TranslationClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *TranslationClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *TranslationClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// TranslateText translates input text and returns translated text.
func (c *TranslationClient) TranslateText(ctx context.Context, req *translatepb.TranslateTextRequest, opts ...gax.CallOption) (*translatepb.TranslateTextResponse, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.TranslateText[0:len(c.CallOptions.TranslateText):len(c.CallOptions.TranslateText)], opts...)
	var resp *translatepb.TranslateTextResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.translationClient.TranslateText(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DetectLanguage detects the language of text within a request.
func (c *TranslationClient) DetectLanguage(ctx context.Context, req *translatepb.DetectLanguageRequest, opts ...gax.CallOption) (*translatepb.DetectLanguageResponse, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.DetectLanguage[0:len(c.CallOptions.DetectLanguage):len(c.CallOptions.DetectLanguage)], opts...)
	var resp *translatepb.DetectLanguageResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.translationClient.DetectLanguage(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetSupportedLanguages returns a list of supported languages for translation.
func (c *TranslationClient) GetSupportedLanguages(ctx context.Context, req *translatepb.GetSupportedLanguagesRequest, opts ...gax.CallOption) (*translatepb.SupportedLanguages, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetSupportedLanguages[0:len(c.CallOptions.GetSupportedLanguages):len(c.CallOptions.GetSupportedLanguages)], opts...)
	var resp *translatepb.SupportedLanguages
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.translationClient.GetSupportedLanguages(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// BatchTranslateText translates a large volume of text in asynchronous batch mode.
// This function provides real-time output as the inputs are being processed.
// If caller cancels a request, the partial results (for an input file, it's
// all or nothing) may still be available on the specified output location.
//
// This call returns immediately and you can
// use google.longrunning.Operation.name to poll the status of the call.
func (c *TranslationClient) BatchTranslateText(ctx context.Context, req *translatepb.BatchTranslateTextRequest, opts ...gax.CallOption) (*BatchTranslateTextOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.BatchTranslateText[0:len(c.CallOptions.BatchTranslateText):len(c.CallOptions.BatchTranslateText)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.translationClient.BatchTranslateText(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &BatchTranslateTextOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// CreateGlossary creates a glossary and returns the long-running operation. Returns
// NOT_FOUND, if the project doesn't exist.
func (c *TranslationClient) CreateGlossary(ctx context.Context, req *translatepb.CreateGlossaryRequest, opts ...gax.CallOption) (*CreateGlossaryOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.CreateGlossary[0:len(c.CallOptions.CreateGlossary):len(c.CallOptions.CreateGlossary)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.translationClient.CreateGlossary(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &CreateGlossaryOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// ListGlossaries lists glossaries in a project. Returns NOT_FOUND, if the project doesn't
// exist.
func (c *TranslationClient) ListGlossaries(ctx context.Context, req *translatepb.ListGlossariesRequest, opts ...gax.CallOption) *GlossaryIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ListGlossaries[0:len(c.CallOptions.ListGlossaries):len(c.CallOptions.ListGlossaries)], opts...)
	it := &GlossaryIterator{}
	req = proto.Clone(req).(*translatepb.ListGlossariesRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*translatepb.Glossary, string, error) {
		var resp *translatepb.ListGlossariesResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.translationClient.ListGlossaries(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Glossaries, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.PageSize)
	it.pageInfo.Token = req.PageToken
	return it
}

// GetGlossary gets a glossary. Returns NOT_FOUND, if the glossary doesn't
// exist.
func (c *TranslationClient) GetGlossary(ctx context.Context, req *translatepb.GetGlossaryRequest, opts ...gax.CallOption) (*translatepb.Glossary, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetGlossary[0:len(c.CallOptions.GetGlossary):len(c.CallOptions.GetGlossary)], opts...)
	var resp *translatepb.Glossary
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.translationClient.GetGlossary(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteGlossary deletes a glossary, or cancels glossary construction
// if the glossary isn't created yet.
// Returns NOT_FOUND, if the glossary doesn't exist.
func (c *TranslationClient) DeleteGlossary(ctx context.Context, req *translatepb.DeleteGlossaryRequest, opts ...gax.CallOption) (*DeleteGlossaryOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.DeleteGlossary[0:len(c.CallOptions.DeleteGlossary):len(c.CallOptions.DeleteGlossary)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.translationClient.DeleteGlossary(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &DeleteGlossaryOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// GlossaryIterator manages a stream of *translatepb.Glossary.
type GlossaryIterator struct {
	items    []*translatepb.Glossary
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*translatepb.Glossary, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *GlossaryIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *GlossaryIterator) Next() (*translatepb.Glossary, error) {
	var item *translatepb.Glossary
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *GlossaryIterator) bufLen() int {
	return len(it.items)
}

func (it *GlossaryIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

// BatchTranslateTextOperation manages a long-running operation from BatchTranslateText.
type BatchTranslateTextOperation struct {
	lro *longrunning.Operation
}

// BatchTranslateTextOperation returns a new BatchTranslateTextOperation from a given name.
// The name must be that of a previously created BatchTranslateTextOperation, possibly from a different process.
func (c *TranslationClient) BatchTranslateTextOperation(name string) *BatchTranslateTextOperation {
	return &BatchTranslateTextOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *BatchTranslateTextOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*translatepb.BatchTranslateResponse, error) {
	var resp translatepb.BatchTranslateResponse
	if err := op.lro.WaitWithInterval(ctx, &resp, 5000*time.Millisecond, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *BatchTranslateTextOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*translatepb.BatchTranslateResponse, error) {
	var resp translatepb.BatchTranslateResponse
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *BatchTranslateTextOperation) Metadata() (*translatepb.BatchTranslateMetadata, error) {
	var meta translatepb.BatchTranslateMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *BatchTranslateTextOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *BatchTranslateTextOperation) Name() string {
	return op.lro.Name()
}

// CreateGlossaryOperation manages a long-running operation from CreateGlossary.
type CreateGlossaryOperation struct {
	lro *longrunning.Operation
}

// CreateGlossaryOperation returns a new CreateGlossaryOperation from a given name.
// The name must be that of a previously created CreateGlossaryOperation, possibly from a different process.
func (c *TranslationClient) CreateGlossaryOperation(name string) *CreateGlossaryOperation {
	return &CreateGlossaryOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *CreateGlossaryOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*translatepb.Glossary, error) {
	var resp translatepb.Glossary
	if err := op.lro.WaitWithInterval(ctx, &resp, 5000*time.Millisecond, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *CreateGlossaryOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*translatepb.Glossary, error) {
	var resp translatepb.Glossary
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *CreateGlossaryOperation) Metadata() (*translatepb.CreateGlossaryMetadata, error) {
	var meta translatepb.CreateGlossaryMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *CreateGlossaryOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *CreateGlossaryOperation) Name() string {
	return op.lro.Name()
}

// DeleteGlossaryOperation manages a long-running operation from DeleteGlossary.
type DeleteGlossaryOperation struct {
	lro *longrunning.Operation
}

// DeleteGlossaryOperation returns a new DeleteGlossaryOperation from a given name.
// The name must be that of a previously created DeleteGlossaryOperation, possibly from a different process.
func (c *TranslationClient) DeleteGlossaryOperation(name string) *DeleteGlossaryOperation {
	return &DeleteGlossaryOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *DeleteGlossaryOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*translatepb.DeleteGlossaryResponse, error) {
	var resp translatepb.DeleteGlossaryResponse
	if err := op.lro.WaitWithInterval(ctx, &resp, 5000*time.Millisecond, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *DeleteGlossaryOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*translatepb.DeleteGlossaryResponse, error) {
	var resp translatepb.DeleteGlossaryResponse
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *DeleteGlossaryOperation) Metadata() (*translatepb.DeleteGlossaryMetadata, error) {
	var meta translatepb.DeleteGlossaryMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *DeleteGlossaryOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *DeleteGlossaryOperation) Name() string {
	return op.lro.Name()
}
