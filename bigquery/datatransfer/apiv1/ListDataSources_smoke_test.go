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

package datatransfer

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/smyte/google-cloud-go/go/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	datatransferpb "google.golang.org/genproto/googleapis/cloud/bigquery/datatransfer/v1"
)

var _ = fmt.Sprintf
var _ = iterator.Done
var _ = strconv.FormatUint
var _ = time.Now

func TestDataTransferServiceSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping smoke test in short mode")
	}
	ctx := context.Background()
	ts := testutil.TokenSource(ctx, DefaultAuthScopes()...)
	if ts == nil {
		t.Skip("Integration tests skipped. See CONTRIBUTING.md for details")
	}

	projectId := testutil.ProjID()
	_ = projectId

	c, err := NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		t.Fatal(err)
	}

	var formattedParent string = fmt.Sprintf("projects/%s/locations/%s", projectId, "us-central1")
	var request = &datatransferpb.ListDataSourcesRequest{
		Parent: formattedParent,
	}

	iter := c.ListDataSources(ctx, request)
	if _, err := iter.Next(); err != nil && err != iterator.Done {
		t.Error(err)
	}
}
