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

package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"

	eventstore "github.com/triggermesh/eventstore/pkg/protob"
)

// EventStore client communicates with an stateful store providing
// global, bridge, and instance access functions.
type EventStore interface {
	Connect(ctx context.Context) error
	Disconnect() error
	Global() Interface
	Bridge(string) Interface
	Instance(string, string) Interface
}

// Interface provides read, write and delete primitives at
// the EventStore
type Interface interface {
	LoadValue(ctx context.Context, key string) ([]byte, error)
	SaveValue(ctx context.Context, key string, value []byte, ttlSec int32) error
	DeleteValue(ctx context.Context, key string) error
}

// client is the default implementation of the stateful
// store client interface.
type client struct {
	// stateful store URI.
	uri string
	// timeout for stateful requests
	timeout time.Duration

	conn     *grpc.ClientConn
	esClient eventstore.EventStoreClient
}

type internalClient struct {
	esClient eventstore.EventStoreClient

	bridge   string
	instance string
}

// New creates an instance of the EventStore client.
func New(uri string, timeout time.Duration) EventStore {
	return &client{
		uri:     uri,
		timeout: timeout,
	}
}

// Connect to the EventStore
func (c *client) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.uri, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("Could not connect to store at %s: %w", c.uri, err)
	}

	c.conn = conn
	c.esClient = eventstore.NewEventStoreClient(conn)

	return nil
}

// Disconnect from the EventStore
func (c *client) Disconnect() error {
	if c.esClient != nil {
		c.esClient = nil
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}
	c.conn = nil

	return nil
}

// Global returns a client that uses the
// brige level to perform storage operations
func (c *client) Global() Interface {
	return &internalClient{
		esClient: c.esClient,
	}
}

// Bridge returns a client that uses the
// brige level to perform storage operations
func (c *client) Bridge(name string) Interface {
	return &internalClient{
		esClient: c.esClient,
		bridge:   name,
	}
}

// Instance returns a client that uses the
// instance level to perform storage operations
func (c *client) Instance(bridge, instance string) Interface {
	return &internalClient{
		esClient: c.esClient,
		bridge:   bridge,
		instance: instance,
	}
}

// LoadValue from EventStore
func (ic *internalClient) LoadValue(ctx context.Context, key string) ([]byte, error) {
	if ic.esClient == nil {
		return nil, errors.New("Event store client is not connected")
	}

	lr := &eventstore.LoadRequest{
		Location: &eventstore.LocationType{
			Key: key,
			Scope: &eventstore.ScopeType{
				Bridge:   ic.bridge,
				Instance: ic.instance,
			}}}

	switch {
	case ic.instance != "":
		lr.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case ic.bridge != "":
		lr.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		lr.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := lr.Validate(); err != nil {
		return nil, err
	}

	r, err := ic.esClient.Load(ctx, lr)
	if err != nil {
		return nil, err
	}

	return r.GetValue(), nil
}

// SaveValue to EventStore
func (ic *internalClient) SaveValue(ctx context.Context, key string, value []byte, ttlSec int32) error {
	if ic.esClient == nil {
		return errors.New("Event store client is not connected")
	}

	sr := &eventstore.SaveRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   ic.bridge,
				Instance: ic.instance,
			},
			Key: key,
		},
		Ttl:   ttlSec,
		Value: value,
	}

	switch {
	case ic.instance != "":
		sr.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case ic.bridge != "":
		sr.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		sr.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := sr.Validate(); err != nil {
		return err
	}

	_, err := ic.esClient.Save(ctx, sr)
	return err
}

// DeleteValue from EventStore
func (ic *internalClient) DeleteValue(ctx context.Context, key string) error {
	if ic.esClient == nil {
		return errors.New("Event store client is not connected")
	}

	dr := &eventstore.DeleteRequest{
		Location: &eventstore.LocationType{
			Key: key,
			Scope: &eventstore.ScopeType{
				Bridge:   ic.bridge,
				Instance: ic.instance,
			}}}

	switch {
	case ic.instance != "":
		dr.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case ic.bridge != "":
		dr.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		dr.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := dr.Validate(); err != nil {
		return err
	}

	_, err := ic.esClient.Delete(ctx, dr)
	return err
}
