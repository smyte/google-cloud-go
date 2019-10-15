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

package texttospeech

import (
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/smyte/google-cloud-go/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var _ = fmt.Sprintf
var _ = iterator.Done
var _ = strconv.FormatUint
var _ = time.Now

func TestTextToSpeechSmoke(t *testing.T) {
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

	var text string = "test"
	var input = &texttospeechpb.SynthesisInput{
		InputSource: &texttospeechpb.SynthesisInput_Text{
			Text: text,
		},
	}
	var languageCode string = "en-US"
	var voice = &texttospeechpb.VoiceSelectionParams{
		LanguageCode: languageCode,
	}
	var audioEncoding texttospeechpb.AudioEncoding = texttospeechpb.AudioEncoding_MP3
	var audioConfig = &texttospeechpb.AudioConfig{
		AudioEncoding: audioEncoding,
	}
	var request = &texttospeechpb.SynthesizeSpeechRequest{
		Input:       input,
		Voice:       voice,
		AudioConfig: audioConfig,
	}

	if _, err := c.SynthesizeSpeech(ctx, request); err != nil {
		t.Error(err)
	}
}
