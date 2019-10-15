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

package dataproc_test

import (
	"context"

	dataproc "github.com/smyte/google-cloud-go/go/dataproc/apiv1beta2"
	"google.golang.org/api/iterator"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1beta2"
)

func ExampleNewAutoscalingPolicyClient() {
	ctx := context.Background()
	c, err := dataproc.NewAutoscalingPolicyClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleAutoscalingPolicyClient_CreateAutoscalingPolicy() {
	ctx := context.Background()
	c, err := dataproc.NewAutoscalingPolicyClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &dataprocpb.CreateAutoscalingPolicyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateAutoscalingPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleAutoscalingPolicyClient_UpdateAutoscalingPolicy() {
	ctx := context.Background()
	c, err := dataproc.NewAutoscalingPolicyClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &dataprocpb.UpdateAutoscalingPolicyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateAutoscalingPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleAutoscalingPolicyClient_GetAutoscalingPolicy() {
	ctx := context.Background()
	c, err := dataproc.NewAutoscalingPolicyClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &dataprocpb.GetAutoscalingPolicyRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetAutoscalingPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleAutoscalingPolicyClient_ListAutoscalingPolicies() {
	ctx := context.Background()
	c, err := dataproc.NewAutoscalingPolicyClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &dataprocpb.ListAutoscalingPoliciesRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListAutoscalingPolicies(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		// TODO: Use resp.
		_ = resp
	}
}

func ExampleAutoscalingPolicyClient_DeleteAutoscalingPolicy() {
	ctx := context.Background()
	c, err := dataproc.NewAutoscalingPolicyClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &dataprocpb.DeleteAutoscalingPolicyRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteAutoscalingPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}
