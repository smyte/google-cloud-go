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

package logadmin_test

import (
	"context"
	"fmt"

	"github.com/smyte/google-cloud-go/logging/logadmin"
	"google.golang.org/api/iterator"
)

func ExampleClient_Sinks() {
	ctx := context.Background()
	client, err := logadmin.NewClient(ctx, "my-project")
	if err != nil {
		// TODO: Handle error.
	}
	it := client.Sinks(ctx)
	_ = it // TODO: iterate using Next or iterator.Pager.
}

func ExampleSinkIterator_Next() {
	ctx := context.Background()
	client, err := logadmin.NewClient(ctx, "my-project")
	if err != nil {
		// TODO: Handle error.
	}
	it := client.Sinks(ctx)
	for {
		sink, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		fmt.Println(sink)
	}
}
