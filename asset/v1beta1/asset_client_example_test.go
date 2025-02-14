// Copyright 2018 Google LLC
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

// AUTO-GENERATED CODE. DO NOT EDIT.

package asset_test

import (
	"context"

	asset "github.com/smyte/google-cloud-go/asset/v1beta1"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1beta1"
)

func ExampleNewClient() {
	ctx := context.Background()
	c, err := asset.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleClient_ExportAssets() {
	ctx := context.Background()
	c, err := asset.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &assetpb.ExportAssetsRequest{
		// TODO: Fill request struct fields.
	}
	op, err := c.ExportAssets(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleClient_BatchGetAssetsHistory() {
	ctx := context.Background()
	c, err := asset.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &assetpb.BatchGetAssetsHistoryRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.BatchGetAssetsHistory(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
