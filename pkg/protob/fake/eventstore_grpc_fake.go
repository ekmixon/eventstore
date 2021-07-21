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
	"errors"

	"github.com/triggermesh/eventstore/pkg/protob"
	"google.golang.org/grpc"
)

// KVStoreClient is a mocked EventStore client
type KVStoreClient interface {
	protob.KVStoreClient

	GetRequests() []Request
}

// Request is a generic placeholder for requests
type Request struct {
	Operation string
	Location  protob.LocationType
	Value     []byte
	TTL       int32
}

// KVStoreOption for customizing the fake client
type KVStoreOption func(KVStoreClient)

// SetMock is the function that mocks save calls
type SetMock func(*protob.SetKVRequest) (*protob.SetKVResponse, error)

// GetMock is the function that mocks save calls
type GetMock func(*protob.GetKVRequest) (*protob.GetKVResponse, error)

// DelMock is the function that mocks save calls
type DelMock func(*protob.DelKVRequest) (*protob.DelKVResponse, error)

type client struct {
	requests []Request
	set      SetMock
	get      GetMock
	del      DelMock
}

// NewEventStoreClientFake creates a fake client
func NewEventStoreClientFake(opts ...KVStoreOption) KVStoreClient {
	c := &client{
		requests: []Request{},
		set:      func(*protob.SetKVRequest) (*protob.SetKVResponse, error) { return nil, nil },
		get:      func(*protob.GetKVRequest) (*protob.GetKVResponse, error) { return nil, nil },
		del:      func(*protob.DelKVRequest) (*protob.DelKVResponse, error) { return nil, nil },
	}

	for _, f := range opts {
		f(c)
	}
	return c
}

// WithSave adds Save function mock to the fake client
func WithSet(f SetMock) KVStoreOption {
	return func(esc KVStoreClient) {
		c := esc.(*client)
		c.set = f
	}
}

// WithLoad adds Load function mock to the fake client
func WithLoad(f GetMock) KVStoreOption {
	return func(esc KVStoreClient) {
		c := esc.(*client)
		c.get = f
	}
}

// WithDelete adds Delete function mock to the fake client
func WithDelete(f DelMock) KVStoreOption {
	return func(esc KVStoreClient) {
		c := esc.(*client)
		c.del = f
	}
}

// GetRequests return the list of requests received at the client
func (c *client) GetRequests() []Request {
	return c.requests
}

// Save variable to storage
func (c *client) Set(ctx context.Context, in *protob.SetKVRequest, opts ...grpc.CallOption) (*protob.SetKVResponse, error) {
	c.requests = append(c.requests, Request{
		Operation: "Set",
		//nolint:govet
		Location: *in.Location,
		Value:    in.Value,
		TTL:      in.Ttl,
	})

	return c.set(in)
}

// Load variable from storage
func (c *client) Get(ctx context.Context, in *protob.GetKVRequest, opts ...grpc.CallOption) (*protob.GetKVResponse, error) {
	c.requests = append(c.requests, Request{
		Operation: "Get",
		//nolint:govet
		Location: *in.Location,
	})

	return c.get(in)
}

// Delete variable from storage
func (c *client) Del(ctx context.Context, in *protob.DelKVRequest, opts ...grpc.CallOption) (*protob.DelKVResponse, error) {
	c.requests = append(c.requests, Request{
		Operation: "Del",
		//nolint:govet
		Location: *in.Location,
	})

	return c.del(in)
}

func (c *client) Incr(ctx context.Context, in *protob.IncrKVRequest, opts ...grpc.CallOption) (*protob.IncrKVResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *client) Decr(ctx context.Context, in *protob.DecrKVRequest, opts ...grpc.CallOption) (*protob.DecrKVResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *client) Lock(ctx context.Context, in *protob.LockRequest, opts ...grpc.CallOption) (*protob.LockResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *client) Unlock(ctx context.Context, in *protob.UnlockRequest, opts ...grpc.CallOption) (*protob.UnlockResponse, error) {
	return nil, errors.New("not implemented")
}
