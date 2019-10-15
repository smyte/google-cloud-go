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

package errorreporting_test

import (
	"context"

	errorreporting "github.com/smyte/google-cloud-go/errorreporting/apiv1beta1"
	clouderrorreportingpb "google.golang.org/genproto/googleapis/devtools/clouderrorreporting/v1beta1"
)

func ExampleNewReportErrorsClient() {
	ctx := context.Background()
	c, err := errorreporting.NewReportErrorsClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleReportErrorsClient_ReportErrorEvent() {
	ctx := context.Background()
	c, err := errorreporting.NewReportErrorsClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &clouderrorreportingpb.ReportErrorEventRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.ReportErrorEvent(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
