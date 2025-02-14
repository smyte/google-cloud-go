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

package scheduler_test

import (
	"context"

	scheduler "github.com/smyte/google-cloud-go/scheduler/apiv1beta1"
	"google.golang.org/api/iterator"
	schedulerpb "google.golang.org/genproto/googleapis/cloud/scheduler/v1beta1"
)

func ExampleNewCloudSchedulerClient() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleCloudSchedulerClient_ListJobs() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.ListJobsRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListJobs(ctx, req)
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

func ExampleCloudSchedulerClient_GetJob() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.GetJobRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetJob(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudSchedulerClient_CreateJob() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.CreateJobRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateJob(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudSchedulerClient_UpdateJob() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.UpdateJobRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateJob(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudSchedulerClient_DeleteJob() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.DeleteJobRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteJob(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}

func ExampleCloudSchedulerClient_PauseJob() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.PauseJobRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.PauseJob(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudSchedulerClient_ResumeJob() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.ResumeJobRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.ResumeJob(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleCloudSchedulerClient_RunJob() {
	ctx := context.Background()
	c, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &schedulerpb.RunJobRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.RunJob(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
