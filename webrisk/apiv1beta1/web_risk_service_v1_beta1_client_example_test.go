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

package webrisk_test

import (
	"context"

	webrisk "github.com/smyte/google-cloud-go/webrisk/apiv1beta1"
	webriskpb "google.golang.org/genproto/googleapis/cloud/webrisk/v1beta1"
)

func ExampleNewWebRiskServiceV1Beta1Client() {
	ctx := context.Background()
	c, err := webrisk.NewWebRiskServiceV1Beta1Client(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleWebRiskServiceV1Beta1Client_ComputeThreatListDiff() {
	ctx := context.Background()
	c, err := webrisk.NewWebRiskServiceV1Beta1Client(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &webriskpb.ComputeThreatListDiffRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.ComputeThreatListDiff(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleWebRiskServiceV1Beta1Client_SearchUris() {
	ctx := context.Background()
	c, err := webrisk.NewWebRiskServiceV1Beta1Client(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &webriskpb.SearchUrisRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.SearchUris(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleWebRiskServiceV1Beta1Client_SearchHashes() {
	ctx := context.Background()
	c, err := webrisk.NewWebRiskServiceV1Beta1Client(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &webriskpb.SearchHashesRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.SearchHashes(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
