/*
Copyright (c) 2020 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"context"

	"github.com/triggermesh/eventstore/pkg/protob"
	"google.golang.org/grpc"
)

// EventStoreClient is a mocked EventStore client
type EventStoreClient interface {
	protob.EventStoreClient

	GetRequests() []Request
}

// Request is a generic placeholder for requests
type Request struct {
	Operation string
	Location  protob.LocationType
	Value     []byte
	TTL       int32
}

// EventStoreOption for customizing the fake client
type EventStoreOption func(EventStoreClient)

// SaveMock is the function that mocks save calls
type SaveMock func(*protob.SaveRequest) (*protob.SaveResponse, error)

// LoadMock is the function that mocks save calls
type LoadMock func(*protob.LoadRequest) (*protob.LoadResponse, error)

// DeleteMock is the function that mocks save calls
type DeleteMock func(*protob.DeleteRequest) (*protob.DeleteResponse, error)

type client struct {
	requests []Request
	save     SaveMock
	load     LoadMock
	delete   DeleteMock
}

// NewEventStoreClientFake creates a fake client
func NewEventStoreClientFake(opts ...EventStoreOption) EventStoreClient {
	c := &client{
		requests: []Request{},
		save:     func(*protob.SaveRequest) (*protob.SaveResponse, error) { return nil, nil },
		load:     func(*protob.LoadRequest) (*protob.LoadResponse, error) { return nil, nil },
		delete:   func(*protob.DeleteRequest) (*protob.DeleteResponse, error) { return nil, nil },
	}

	for _, f := range opts {
		f(c)
	}
	return c
}

// WithSave adds Save function mock to the fake client
func WithSave(f SaveMock) EventStoreOption {
	return func(esc EventStoreClient) {
		c := esc.(*client)
		c.save = f
	}
}

// WithLoad adds Load function mock to the fake client
func WithLoad(f LoadMock) EventStoreOption {
	return func(esc EventStoreClient) {
		c := esc.(*client)
		c.load = f
	}
}

// WithDelete adds Delete function mock to the fake client
func WithDelete(f DeleteMock) EventStoreOption {
	return func(esc EventStoreClient) {
		c := esc.(*client)
		c.delete = f
	}
}

// GetRequests return the list of requests received at the client
func (c *client) GetRequests() []Request {
	return c.requests
}

// Save variable to storage
func (c *client) Save(ctx context.Context, in *protob.SaveRequest, opts ...grpc.CallOption) (*protob.SaveResponse, error) {
	c.requests = append(c.requests, Request{
		Operation: "Save",
		//nolint:govet
		Location: *in.Location,
		Value:    in.Value,
		TTL:      in.Ttl,
	})

	return c.save(in)
}

// Load variable from storage
func (c *client) Load(ctx context.Context, in *protob.LoadRequest, opts ...grpc.CallOption) (*protob.LoadResponse, error) {
	c.requests = append(c.requests, Request{
		Operation: "Load",
		//nolint:govet
		Location: *in.Location,
	})

	return c.load(in)
}

// Delete variable from storage
func (c *client) Delete(ctx context.Context, in *protob.DeleteRequest, opts ...grpc.CallOption) (*protob.DeleteResponse, error) {
	c.requests = append(c.requests, Request{
		Operation: "Delete",
		//nolint:govet
		Location: *in.Location,
	})

	return c.delete(in)
}
