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

package talent_test

import (
	"context"

	talent "github.com/smyte/google-cloud-go/go/talent/apiv4beta1"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

func ExampleNewCompletionClient() {
	ctx := context.Background()
	c, err := talent.NewCompletionClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleCompletionClient_CompleteQuery() {
	ctx := context.Background()
	c, err := talent.NewCompletionClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &talentpb.CompleteQueryRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CompleteQuery(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
