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

package pubsub_test

import (
	"context"

	pubsub "github.com/smyte/google-cloud-go/pubsub/apiv1"
	"google.golang.org/api/iterator"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

func ExampleNewPublisherClient() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExamplePublisherClient_CreateTopic() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &pubsubpb.Topic{
		// TODO: Fill request struct fields.
	}
	resp, err := c.CreateTopic(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExamplePublisherClient_UpdateTopic() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &pubsubpb.UpdateTopicRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.UpdateTopic(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExamplePublisherClient_Publish() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &pubsubpb.PublishRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.Publish(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExamplePublisherClient_GetTopic() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &pubsubpb.GetTopicRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetTopic(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExamplePublisherClient_ListTopics() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &pubsubpb.ListTopicsRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListTopics(ctx, req)
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

func ExamplePublisherClient_ListTopicSubscriptions() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &pubsubpb.ListTopicSubscriptionsRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListTopicSubscriptions(ctx, req)
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

func ExamplePublisherClient_DeleteTopic() {
	ctx := context.Background()
	c, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &pubsubpb.DeleteTopicRequest{
		// TODO: Fill request struct fields.
	}
	err = c.DeleteTopic(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
}
