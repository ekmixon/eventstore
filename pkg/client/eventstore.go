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
	KV() KeyValue
	Map() Map
	Queue() Queue
	Sync() Sync
}

type Sync interface {
	Lock(ctx context.Context, key string, timeout int32) (string, error)
	Unlock(ctx context.Context, key string, unlock string) error
}

// KeyValue is the key value interface for storage.
type KeyValue interface {
	Set(ctx context.Context, key string, value []byte, ttlSec int32) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string, value int32) error
	Decr(ctx context.Context, key string, value int32) error
}

// MapInterface is the map structure interface for storage.
type Map interface {
	New(ctx context.Context, key string, ttlSec int32) error
	Fields(key string) MapFields
	Del(ctx context.Context, key string) error
}

type MapFields interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string, value int32) error
	Decr(ctx context.Context, key string, value int32) error

	All(ctx context.Context) (map[string][]byte, error)
	Len(ctx context.Context) (int, error)
}

// Queue is a minimal FIFO interface.
type Queue interface {
	New(ctx context.Context, key string, ttlSec int32) error
	Items(key string) QueueItems
	Del(ctx context.Context, key string) error
}
type QueueItems interface {
	Push(ctx context.Context, value []byte) error
	Index(ctx context.Context, index int32) ([]byte, error)
	Pop(ctx context.Context) ([]byte, error)
	Peek(ctx context.Context) ([]byte, error)

	All(ctx context.Context) ([][]byte, error)
	Len(ctx context.Context) (int, error)
}

// client is the default implementation of the stateful
// store client interface.
type client struct {
	// stateful store URI.
	uri string
	// timeout for stateful requests
	timeout  time.Duration
	conn     *grpc.ClientConn
	services *services
}

type services struct {
	kvc    eventstore.KVClient
	mapc   eventstore.MapClient
	queuec eventstore.QueueClient
	syncc  eventstore.SyncClient
}

type internalClient struct {
	svc *services

	bridge   string
	instance string
}

func (s *internalClient) KV() KeyValue {
	return &internalKV{s}
}

func (s *internalClient) Map() Map {
	return &internalMap{s}
}

func (s *internalClient) Queue() Queue {
	return &internalQueue{s}
}

func (s *internalClient) Sync() Sync {
	// TODO add Queue
	return &internalSync{s}
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
		return fmt.Errorf("could not connect to store at %s: %w", c.uri, err)
	}

	c.services = &services{
		kvc:    eventstore.NewKVClient(conn),
		mapc:   eventstore.NewMapClient(conn),
		queuec: eventstore.NewQueueClient(conn),
		syncc:  eventstore.NewSyncClient(conn),
	}

	return nil
}

// Disconnect from the EventStore
func (c *client) Disconnect() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}

	c.services.kvc = nil
	c.services.mapc = nil
	c.services.queuec = nil
	c.services.syncc = nil
	c.conn = nil

	return nil
}

// Global returns a client that uses the
// brige level to perform storage operations
func (c *client) Global() Interface {
	return &internalClient{
		svc: c.services,
	}
}

// Bridge returns a client that uses the
// brige level to perform storage operations
func (c *client) Bridge(name string) Interface {
	return &internalClient{
		svc:    c.services,
		bridge: name,
	}
}

// Instance returns a client that uses the
// instance level to perform storage operations
func (c *client) Instance(bridge, instance string) Interface {
	return &internalClient{
		svc:      c.services,
		bridge:   bridge,
		instance: instance,
	}
}
