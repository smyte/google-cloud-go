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

package translate_test

import (
	"context"

	translate "github.com/smyte/google-cloud-go/translate/apiv3"
	"google.golang.org/api/iterator"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

func ExampleNewTranslationClient() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use client.
	_ = c
}

func ExampleTranslationClient_TranslateText() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.TranslateTextRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.TranslateText(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleTranslationClient_DetectLanguage() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.DetectLanguageRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.DetectLanguage(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleTranslationClient_GetSupportedLanguages() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.GetSupportedLanguagesRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetSupportedLanguages(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleTranslationClient_BatchTranslateText() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.BatchTranslateTextRequest{
		// TODO: Fill request struct fields.
	}
	op, err := c.BatchTranslateText(ctx, req)
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

func ExampleTranslationClient_CreateGlossary() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.CreateGlossaryRequest{
		// TODO: Fill request struct fields.
	}
	op, err := c.CreateGlossary(ctx, req)
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

func ExampleTranslationClient_ListGlossaries() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.ListGlossariesRequest{
		// TODO: Fill request struct fields.
	}
	it := c.ListGlossaries(ctx, req)
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

func ExampleTranslationClient_GetGlossary() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.GetGlossaryRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := c.GetGlossary(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}

func ExampleTranslationClient_DeleteGlossary() {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	req := &translatepb.DeleteGlossaryRequest{
		// TODO: Fill request struct fields.
	}
	op, err := c.DeleteGlossary(ctx, req)
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
