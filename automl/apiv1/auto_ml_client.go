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

package automl

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/smyte/google-cloud-go/longrunning"
	lroauto "github.com/smyte/google-cloud-go/longrunning/autogen"
	"github.com/golang/protobuf/proto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// CallOptions contains the retry settings for each method of Client.
type CallOptions struct {
	CreateDataset        []gax.CallOption
	UpdateDataset        []gax.CallOption
	GetDataset           []gax.CallOption
	ListDatasets         []gax.CallOption
	DeleteDataset        []gax.CallOption
	ImportData           []gax.CallOption
	ExportData           []gax.CallOption
	CreateModel          []gax.CallOption
	GetModel             []gax.CallOption
	UpdateModel          []gax.CallOption
	ListModels           []gax.CallOption
	DeleteModel          []gax.CallOption
	GetModelEvaluation   []gax.CallOption
	ListModelEvaluations []gax.CallOption
}

func defaultClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("automl.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultCallOptions() *CallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.3,
				})
			}),
		},
	}
	return &CallOptions{
		CreateDataset:        retry[[2]string{"default", "non_idempotent"}],
		UpdateDataset:        retry[[2]string{"default", "non_idempotent"}],
		GetDataset:           retry[[2]string{"default", "idempotent"}],
		ListDatasets:         retry[[2]string{"default", "idempotent"}],
		DeleteDataset:        retry[[2]string{"default", "idempotent"}],
		ImportData:           retry[[2]string{"default", "non_idempotent"}],
		ExportData:           retry[[2]string{"default", "non_idempotent"}],
		CreateModel:          retry[[2]string{"default", "non_idempotent"}],
		GetModel:             retry[[2]string{"default", "idempotent"}],
		UpdateModel:          retry[[2]string{"default", "non_idempotent"}],
		ListModels:           retry[[2]string{"default", "idempotent"}],
		DeleteModel:          retry[[2]string{"default", "idempotent"}],
		GetModelEvaluation:   retry[[2]string{"default", "idempotent"}],
		ListModelEvaluations: retry[[2]string{"default", "idempotent"}],
	}
}

// Client is a client for interacting with Cloud AutoML API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type Client struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	client automlpb.AutoMlClient

	// LROClient is used internally to handle longrunning operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient *lroauto.OperationsClient

	// The call options for this service.
	CallOptions *CallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewClient creates a new auto ml client.
//
// AutoML Server API.
//
// The resource names are assigned by the server.
// The server never reuses names that it has created after the resources with
// those names are deleted.
//
// An ID of a resource is the last element of the item's resource name. For
// projects/{project_id}/locations/{location_id}/datasets/{dataset_id}, then
// the id for the item is {dataset_id}.
//
// Currently the only supported location_id is "us-central1".
//
// On any input that is documented to expect a string parameter in
// snake_case or kebab-case, either of those cases is accepted.
func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn:        conn,
		CallOptions: defaultCallOptions(),

		client: automlpb.NewAutoMlClient(conn),
	}
	c.setGoogleClientInfo()

	c.LROClient, err = lroauto.NewOperationsClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		// This error "should not happen", since we are just reusing old connection
		// and never actually need to dial.
		// If this does happen, we could leak conn. However, we cannot close conn:
		// If the user invoked the function with option.WithGRPCConn,
		// we would close a connection that's still in use.
		// TODO(pongad): investigate error conditions.
		return nil, err
	}
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *Client) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *Client) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *Client) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// CreateDataset creates a dataset.
func (c *Client) CreateDataset(ctx context.Context, req *automlpb.CreateDatasetRequest, opts ...gax.CallOption) (*CreateDatasetOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.CreateDataset[0:len(c.CallOptions.CreateDataset):len(c.CallOptions.CreateDataset)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.CreateDataset(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &CreateDatasetOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// UpdateDataset updates a dataset.
func (c *Client) UpdateDataset(ctx context.Context, req *automlpb.UpdateDatasetRequest, opts ...gax.CallOption) (*automlpb.Dataset, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "dataset.name", url.QueryEscape(req.GetDataset().GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.UpdateDataset[0:len(c.CallOptions.UpdateDataset):len(c.CallOptions.UpdateDataset)], opts...)
	var resp *automlpb.Dataset
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.UpdateDataset(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetDataset gets a dataset.
func (c *Client) GetDataset(ctx context.Context, req *automlpb.GetDatasetRequest, opts ...gax.CallOption) (*automlpb.Dataset, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetDataset[0:len(c.CallOptions.GetDataset):len(c.CallOptions.GetDataset)], opts...)
	var resp *automlpb.Dataset
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.GetDataset(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListDatasets lists datasets in a project.
func (c *Client) ListDatasets(ctx context.Context, req *automlpb.ListDatasetsRequest, opts ...gax.CallOption) *DatasetIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ListDatasets[0:len(c.CallOptions.ListDatasets):len(c.CallOptions.ListDatasets)], opts...)
	it := &DatasetIterator{}
	req = proto.Clone(req).(*automlpb.ListDatasetsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*automlpb.Dataset, string, error) {
		var resp *automlpb.ListDatasetsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.client.ListDatasets(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Datasets, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.PageSize)
	it.pageInfo.Token = req.PageToken
	return it
}

// DeleteDataset deletes a dataset and all of its contents.
// Returns empty response in the
// [response][google.longrunning.Operation.response] field when it completes,
// and delete_details in the
// [metadata][google.longrunning.Operation.metadata] field.
func (c *Client) DeleteDataset(ctx context.Context, req *automlpb.DeleteDatasetRequest, opts ...gax.CallOption) (*DeleteDatasetOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.DeleteDataset[0:len(c.CallOptions.DeleteDataset):len(c.CallOptions.DeleteDataset)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.DeleteDataset(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &DeleteDatasetOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// ImportData imports data into a dataset.
func (c *Client) ImportData(ctx context.Context, req *automlpb.ImportDataRequest, opts ...gax.CallOption) (*ImportDataOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ImportData[0:len(c.CallOptions.ImportData):len(c.CallOptions.ImportData)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.ImportData(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &ImportDataOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// ExportData exports dataset's data to the provided output location.
// Returns an empty response in the
// [response][google.longrunning.Operation.response] field when it completes.
func (c *Client) ExportData(ctx context.Context, req *automlpb.ExportDataRequest, opts ...gax.CallOption) (*ExportDataOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ExportData[0:len(c.CallOptions.ExportData):len(c.CallOptions.ExportData)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.ExportData(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &ExportDataOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// CreateModel creates a model.
// Returns a Model in the [response][google.longrunning.Operation.response]
// field when it completes.
// When you create a model, several model evaluations are created for it:
// a global evaluation, and one evaluation for each annotation spec.
func (c *Client) CreateModel(ctx context.Context, req *automlpb.CreateModelRequest, opts ...gax.CallOption) (*CreateModelOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.CreateModel[0:len(c.CallOptions.CreateModel):len(c.CallOptions.CreateModel)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.CreateModel(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &CreateModelOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// GetModel gets a model.
func (c *Client) GetModel(ctx context.Context, req *automlpb.GetModelRequest, opts ...gax.CallOption) (*automlpb.Model, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetModel[0:len(c.CallOptions.GetModel):len(c.CallOptions.GetModel)], opts...)
	var resp *automlpb.Model
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.GetModel(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateModel updates a model.
func (c *Client) UpdateModel(ctx context.Context, req *automlpb.UpdateModelRequest, opts ...gax.CallOption) (*automlpb.Model, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "model.name", url.QueryEscape(req.GetModel().GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.UpdateModel[0:len(c.CallOptions.UpdateModel):len(c.CallOptions.UpdateModel)], opts...)
	var resp *automlpb.Model
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.UpdateModel(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListModels lists models.
func (c *Client) ListModels(ctx context.Context, req *automlpb.ListModelsRequest, opts ...gax.CallOption) *ModelIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ListModels[0:len(c.CallOptions.ListModels):len(c.CallOptions.ListModels)], opts...)
	it := &ModelIterator{}
	req = proto.Clone(req).(*automlpb.ListModelsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*automlpb.Model, string, error) {
		var resp *automlpb.ListModelsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.client.ListModels(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Model, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.PageSize)
	it.pageInfo.Token = req.PageToken
	return it
}

// DeleteModel deletes a model.
// Returns google.protobuf.Empty in the
// [response][google.longrunning.Operation.response] field when it completes,
// and delete_details in the
// [metadata][google.longrunning.Operation.metadata] field.
func (c *Client) DeleteModel(ctx context.Context, req *automlpb.DeleteModelRequest, opts ...gax.CallOption) (*DeleteModelOperation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.DeleteModel[0:len(c.CallOptions.DeleteModel):len(c.CallOptions.DeleteModel)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.DeleteModel(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &DeleteModelOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, resp),
	}, nil
}

// GetModelEvaluation gets a model evaluation.
func (c *Client) GetModelEvaluation(ctx context.Context, req *automlpb.GetModelEvaluationRequest, opts ...gax.CallOption) (*automlpb.ModelEvaluation, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetModelEvaluation[0:len(c.CallOptions.GetModelEvaluation):len(c.CallOptions.GetModelEvaluation)], opts...)
	var resp *automlpb.ModelEvaluation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.GetModelEvaluation(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListModelEvaluations lists model evaluations.
func (c *Client) ListModelEvaluations(ctx context.Context, req *automlpb.ListModelEvaluationsRequest, opts ...gax.CallOption) *ModelEvaluationIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ListModelEvaluations[0:len(c.CallOptions.ListModelEvaluations):len(c.CallOptions.ListModelEvaluations)], opts...)
	it := &ModelEvaluationIterator{}
	req = proto.Clone(req).(*automlpb.ListModelEvaluationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*automlpb.ModelEvaluation, string, error) {
		var resp *automlpb.ListModelEvaluationsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.client.ListModelEvaluations(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.ModelEvaluation, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.PageSize)
	it.pageInfo.Token = req.PageToken
	return it
}

// DatasetIterator manages a stream of *automlpb.Dataset.
type DatasetIterator struct {
	items    []*automlpb.Dataset
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*automlpb.Dataset, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *DatasetIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *DatasetIterator) Next() (*automlpb.Dataset, error) {
	var item *automlpb.Dataset
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *DatasetIterator) bufLen() int {
	return len(it.items)
}

func (it *DatasetIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

// ModelEvaluationIterator manages a stream of *automlpb.ModelEvaluation.
type ModelEvaluationIterator struct {
	items    []*automlpb.ModelEvaluation
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*automlpb.ModelEvaluation, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *ModelEvaluationIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *ModelEvaluationIterator) Next() (*automlpb.ModelEvaluation, error) {
	var item *automlpb.ModelEvaluation
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *ModelEvaluationIterator) bufLen() int {
	return len(it.items)
}

func (it *ModelEvaluationIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

// ModelIterator manages a stream of *automlpb.Model.
type ModelIterator struct {
	items    []*automlpb.Model
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*automlpb.Model, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *ModelIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *ModelIterator) Next() (*automlpb.Model, error) {
	var item *automlpb.Model
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *ModelIterator) bufLen() int {
	return len(it.items)
}

func (it *ModelIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

// CreateDatasetOperation manages a long-running operation from CreateDataset.
type CreateDatasetOperation struct {
	lro *longrunning.Operation
}

// CreateDatasetOperation returns a new CreateDatasetOperation from a given name.
// The name must be that of a previously created CreateDatasetOperation, possibly from a different process.
func (c *Client) CreateDatasetOperation(name string) *CreateDatasetOperation {
	return &CreateDatasetOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *CreateDatasetOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*automlpb.Dataset, error) {
	var resp automlpb.Dataset
	if err := op.lro.WaitWithInterval(ctx, &resp, 5000*time.Millisecond, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *CreateDatasetOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*automlpb.Dataset, error) {
	var resp automlpb.Dataset
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *CreateDatasetOperation) Metadata() (*automlpb.OperationMetadata, error) {
	var meta automlpb.OperationMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *CreateDatasetOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *CreateDatasetOperation) Name() string {
	return op.lro.Name()
}

// CreateModelOperation manages a long-running operation from CreateModel.
type CreateModelOperation struct {
	lro *longrunning.Operation
}

// CreateModelOperation returns a new CreateModelOperation from a given name.
// The name must be that of a previously created CreateModelOperation, possibly from a different process.
func (c *Client) CreateModelOperation(name string) *CreateModelOperation {
	return &CreateModelOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *CreateModelOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*automlpb.Model, error) {
	var resp automlpb.Model
	if err := op.lro.WaitWithInterval(ctx, &resp, 5000*time.Millisecond, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *CreateModelOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*automlpb.Model, error) {
	var resp automlpb.Model
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *CreateModelOperation) Metadata() (*automlpb.OperationMetadata, error) {
	var meta automlpb.OperationMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *CreateModelOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *CreateModelOperation) Name() string {
	return op.lro.Name()
}

// DeleteDatasetOperation manages a long-running operation from DeleteDataset.
type DeleteDatasetOperation struct {
	lro *longrunning.Operation
}

// DeleteDatasetOperation returns a new DeleteDatasetOperation from a given name.
// The name must be that of a previously created DeleteDatasetOperation, possibly from a different process.
func (c *Client) DeleteDatasetOperation(name string) *DeleteDatasetOperation {
	return &DeleteDatasetOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning any error encountered.
//
// See documentation of Poll for error-handling information.
func (op *DeleteDatasetOperation) Wait(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.WaitWithInterval(ctx, nil, 5000*time.Millisecond, opts...)
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully, op.Done will return true.
func (op *DeleteDatasetOperation) Poll(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.Poll(ctx, nil, opts...)
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *DeleteDatasetOperation) Metadata() (*automlpb.OperationMetadata, error) {
	var meta automlpb.OperationMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *DeleteDatasetOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *DeleteDatasetOperation) Name() string {
	return op.lro.Name()
}

// DeleteModelOperation manages a long-running operation from DeleteModel.
type DeleteModelOperation struct {
	lro *longrunning.Operation
}

// DeleteModelOperation returns a new DeleteModelOperation from a given name.
// The name must be that of a previously created DeleteModelOperation, possibly from a different process.
func (c *Client) DeleteModelOperation(name string) *DeleteModelOperation {
	return &DeleteModelOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning any error encountered.
//
// See documentation of Poll for error-handling information.
func (op *DeleteModelOperation) Wait(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.WaitWithInterval(ctx, nil, 5000*time.Millisecond, opts...)
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully, op.Done will return true.
func (op *DeleteModelOperation) Poll(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.Poll(ctx, nil, opts...)
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *DeleteModelOperation) Metadata() (*automlpb.OperationMetadata, error) {
	var meta automlpb.OperationMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *DeleteModelOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *DeleteModelOperation) Name() string {
	return op.lro.Name()
}

// ExportDataOperation manages a long-running operation from ExportData.
type ExportDataOperation struct {
	lro *longrunning.Operation
}

// ExportDataOperation returns a new ExportDataOperation from a given name.
// The name must be that of a previously created ExportDataOperation, possibly from a different process.
func (c *Client) ExportDataOperation(name string) *ExportDataOperation {
	return &ExportDataOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning any error encountered.
//
// See documentation of Poll for error-handling information.
func (op *ExportDataOperation) Wait(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.WaitWithInterval(ctx, nil, 5000*time.Millisecond, opts...)
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully, op.Done will return true.
func (op *ExportDataOperation) Poll(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.Poll(ctx, nil, opts...)
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *ExportDataOperation) Metadata() (*automlpb.OperationMetadata, error) {
	var meta automlpb.OperationMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *ExportDataOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *ExportDataOperation) Name() string {
	return op.lro.Name()
}

// ImportDataOperation manages a long-running operation from ImportData.
type ImportDataOperation struct {
	lro *longrunning.Operation
}

// ImportDataOperation returns a new ImportDataOperation from a given name.
// The name must be that of a previously created ImportDataOperation, possibly from a different process.
func (c *Client) ImportDataOperation(name string) *ImportDataOperation {
	return &ImportDataOperation{
		lro: longrunning.InternalNewOperation(c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning any error encountered.
//
// See documentation of Poll for error-handling information.
func (op *ImportDataOperation) Wait(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.WaitWithInterval(ctx, nil, 5000*time.Millisecond, opts...)
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully, op.Done will return true.
func (op *ImportDataOperation) Poll(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.Poll(ctx, nil, opts...)
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *ImportDataOperation) Metadata() (*automlpb.OperationMetadata, error) {
	var meta automlpb.OperationMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *ImportDataOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *ImportDataOperation) Name() string {
	return op.lro.Name()
}
