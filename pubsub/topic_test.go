// Copyright 2016 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pubsub

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smyte/google-cloud-go/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/support/bundler"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkTopicListing(t *testing.T, c *Client, want []string) {
	topics, err := slurpTopics(c.Topics(context.Background()))
	if err != nil {
		t.Fatalf("error listing topics: %v", err)
	}
	var got []string
	for _, topic := range topics {
		got = append(got, topic.ID())
	}
	if !testutil.Equal(got, want) {
		t.Errorf("topic list: got: %v, want: %v", got, want)
	}
}

// All returns the remaining topics from this iterator.
func slurpTopics(it *TopicIterator) ([]*Topic, error) {
	var topics []*Topic
	for {
		switch topic, err := it.Next(); err {
		case nil:
			topics = append(topics, topic)
		case iterator.Done:
			return topics, nil
		default:
			return nil, err
		}
	}
}

func TestTopicID(t *testing.T) {
	const id = "id"
	c, srv := newFake(t)
	defer c.Close()
	defer srv.Close()

	s := c.Topic(id)
	if got, want := s.ID(), id; got != want {
		t.Errorf("Topic.ID() = %q; want %q", got, want)
	}
}

func TestCreateTopicWithConfig(t *testing.T) {
	c, srv := newFake(t)
	defer c.Close()
	defer srv.Close()

	id := "test-topic"
	want := TopicConfig{
		Labels: map[string]string{"label": "value"},
		MessageStoragePolicy: MessageStoragePolicy{
			AllowedPersistenceRegions: []string{"us-east1"},
		},
		KMSKeyName: "projects/P/locations/L/keyRings/R/cryptoKeys/K",
	}

	topic := mustCreateTopicWithConfig(t, c, id, &want)
	got, err := topic.Config(context.Background())
	if err != nil {
		t.Fatalf("error getting topic config: %v", err)
	}

	if !testutil.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestListTopics(t *testing.T) {
	c, srv := newFake(t)
	defer c.Close()
	defer srv.Close()

	var ids []string
	for i := 1; i <= 4; i++ {
		id := fmt.Sprintf("t%d", i)
		ids = append(ids, id)
		mustCreateTopic(t, c, id)
	}
	checkTopicListing(t, c, ids)
}

func TestListCompletelyEmptyTopics(t *testing.T) {
	c, srv := newFake(t)
	defer c.Close()
	defer srv.Close()

	checkTopicListing(t, c, nil)
}

func TestStopPublishOrder(t *testing.T) {
	// Check that Stop doesn't panic if called before Publish.
	// Also that Publish after Stop returns the right error.
	ctx := context.Background()
	c := &Client{projectID: "projid"}
	topic := c.Topic("t")
	topic.Stop()
	r := topic.Publish(ctx, &Message{})
	_, err := r.Get(ctx)
	if err != errTopicStopped {
		t.Errorf("got %v, want errTopicStopped", err)
	}
}

func TestPublishTimeout(t *testing.T) {
	ctx := context.Background()
	serv, err := testutil.NewServer()
	if err != nil {
		t.Fatal(err)
	}
	pubsubpb.RegisterPublisherServer(serv.Gsrv, &alwaysFailPublish{})
	serv.Start()
	conn, err := grpc.Dial(serv.Addr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	opts := withGRPCHeadersAssertion(t, option.WithGRPCConn(conn))
	c, err := NewClient(ctx, "projectID", opts...)
	if err != nil {
		t.Fatal(err)
	}
	topic := c.Topic("t")
	topic.PublishSettings.Timeout = 3 * time.Second
	r := topic.Publish(ctx, &Message{})
	defer topic.Stop()
	select {
	case <-r.Ready():
		_, err = r.Get(ctx)
		if err != context.DeadlineExceeded {
			t.Fatalf("got %v, want context.DeadlineExceeded", err)
		}
	case <-time.After(2 * topic.PublishSettings.Timeout):
		t.Fatal("timed out")
	}
}

func TestPublishBufferedByteLimit(t *testing.T) {
	ctx := context.Background()
	client, srv := newFake(t)
	defer client.Close()
	defer srv.Close()

	topic := mustCreateTopic(t, client, "topic-small-buffered-byte-limit")
	defer topic.Stop()

	// Test setting BufferedByteLimit to small number of bytes that should fail.
	topic.PublishSettings.BufferedByteLimit = 100

	const messageSizeBytes = 1000

	msg := &Message{Data: bytes.Repeat([]byte{'A'}, int(messageSizeBytes))}
	res := topic.Publish(ctx, msg)

	_, err := res.Get(ctx)
	if err != bundler.ErrOverflow {
		t.Errorf("got %v, want ErrOverflow", err)
	}
}

func TestUpdateTopic_Label(t *testing.T) {
	ctx := context.Background()
	client, srv := newFake(t)
	defer client.Close()
	defer srv.Close()

	topic := mustCreateTopic(t, client, "T")
	config, err := topic.Config(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := TopicConfig{}
	if !testutil.Equal(config, want) {
		t.Errorf("got %+v, want %+v", config, want)
	}

	// replace labels
	labels := map[string]string{"label": "value"}
	config2, err := topic.Update(ctx, TopicConfigToUpdate{
		Labels: labels,
	})
	if err != nil {
		t.Fatal(err)
	}
	want = TopicConfig{
		Labels: labels,
	}
	if !testutil.Equal(config2, want) {
		t.Errorf("got %+v, want %+v", config2, want)
	}

	// delete all labels
	labels = map[string]string{}
	config3, err := topic.Update(ctx, TopicConfigToUpdate{Labels: labels})
	if err != nil {
		t.Fatal(err)
	}
	want.Labels = nil
	if !testutil.Equal(config3, want) {
		t.Errorf("got %+v, want %+v", config3, want)
	}
}

func TestUpdateTopic_MessageStoragePolicy(t *testing.T) {
	ctx := context.Background()
	client, srv := newFake(t)
	defer client.Close()
	defer srv.Close()

	topic := mustCreateTopic(t, client, "T")
	config, err := topic.Config(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := TopicConfig{}
	if !testutil.Equal(config, want) {
		t.Errorf("\ngot  %+v\nwant %+v", config, want)
	}

	// Update message storage policy.
	msp := &MessageStoragePolicy{
		AllowedPersistenceRegions: []string{"us-east1"},
	}
	config2, err := topic.Update(ctx, TopicConfigToUpdate{MessageStoragePolicy: msp})
	if err != nil {
		t.Fatal(err)
	}
	want.MessageStoragePolicy = MessageStoragePolicy{
		AllowedPersistenceRegions: []string{"us-east1"},
	}
	if !testutil.Equal(config2, want) {
		t.Errorf("\ngot  %+v\nwant %+v", config2, want)
	}
}

type alwaysFailPublish struct {
	pubsubpb.PublisherServer
}

func (s *alwaysFailPublish) Publish(ctx context.Context, req *pubsubpb.PublishRequest) (*pubsubpb.PublishResponse, error) {
	return nil, status.Errorf(codes.Unavailable, "try again")
}

func mustCreateTopic(t *testing.T, c *Client, id string) *Topic {
	topic, err := c.CreateTopic(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	return topic
}

func mustCreateTopicWithConfig(t *testing.T, c *Client, id string, tc *TopicConfig) *Topic {
	if tc == nil {
		return mustCreateTopic(t, c, id)
	}
	topic, err := c.CreateTopicWithConfig(context.Background(), id, tc)
	if err != nil {
		t.Fatal(err)
	}
	return topic
}
