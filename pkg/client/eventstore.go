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
	KV() KeyValue
	Map() Map
	Queue() Queue
}

type Lockable interface {
	Lock(ctx context.Context, key string, timeout int32) (string, error)
	Unlock(ctx context.Context, key string, unlock string) error
}

// KeyValue is the key value interface for storage.
type KeyValue interface {
	Set(ctx context.Context, key string, value []byte, ttlSec int32) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string, value int) error
	Decr(ctx context.Context, key string, value int) error

	Lockable
}

// MapInterface is the map structure interface for storage.
type Map interface {
	New(ctx context.Context, key string, ttlSec int32) error
	Get(ctx context.Context, key string) MapFields
	Del(ctx context.Context, key string, opts ...grpc.CallOption) error

	Lockable
}

type MapFields interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string, value int) error
	Decr(ctx context.Context, key string, value int) error

	GetAll(ctx context.Context) (map[string][]byte, error)
	Len(ctx context.Context) (int, error)
}

// Queue is a minimal FIFO interface.
type Queue interface {
	New(ctx context.Context, key string, ttlSec int32) error
	Get(ctx context.Context, key string) QueueItems
	Del(ctx context.Context, key string, opts ...grpc.CallOption) error

	Lockable
}
type QueueItems interface {
	Push(ctx context.Context, value []byte) error
	Index(ctx context.Context, index int) ([]byte, error)
	Pop(ctx context.Context) ([]byte, error)
	Peek(ctx context.Context) ([]byte, error)

	GetAll(ctx context.Context) (map[string][]byte, error)
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
}

type internalClient struct {
	svc *services

	bridge   string
	instance string
}

type internalKV struct {
	*internalClient
}

func (i *internalKV) Set(ctx context.Context, key string, value []byte, ttlSec int32) error {
	if i.svc.kvc == nil {
		return errors.New("Event store client is not connected")
	}

	sr := &eventstore.SetKVRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: key,
		},
		Ttl:   ttlSec,
		Value: value,
	}

	switch {
	case i.instance != "":
		sr.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		sr.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		sr.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := sr.Validate(); err != nil {
		return err
	}

	_, err := i.svc.kvc.Set(ctx, sr)
	return err
}

// LoadValue from EventStore
func (i *internalKV) Get(ctx context.Context, key string) ([]byte, error) {
	if i.svc.kvc == nil {
		return nil, errors.New("Event store client is not connected")
	}

	gr := &eventstore.GetKVRequest{
		Location: &eventstore.LocationType{
			Key: key,
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			}}}

	switch {
	case i.instance != "":
		gr.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		gr.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		gr.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := gr.Validate(); err != nil {
		return nil, err
	}

	r, err := i.svc.kvc.Get(ctx, gr)
	if err != nil {
		return nil, err
	}

	return r.GetValue(), nil
}

// DeleteValue from EventStore
func (i *internalKV) Del(ctx context.Context, key string) error {
	if i.svc.kvc == nil {
		return errors.New("Event store client is not connected")
	}

	dr := &eventstore.DelKVRequest{
		Location: &eventstore.LocationType{
			Key: key,
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			}}}

	switch {
	case i.instance != "":
		dr.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		dr.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		dr.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := dr.Validate(); err != nil {
		return err
	}

	_, err := i.svc.kvc.Del(ctx, dr)
	return err
}

//
// Incr(ctx context.Context, key string, value int) error
// Decr(ctx context.Context, key string, value int) error

// Lockable

func (s *internalClient) KV() KeyValue {
	return &internalKV{s}
}

func (s *internalClient) Map() Map {

}

func (s *internalClient) Queue() Queue {

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

	c.services = &services{
		kvc:    eventstore.NewKVClient(conn),
		mapc:   eventstore.NewMapClient(conn),
		queuec: eventstore.NewQueueClient(conn),
	}
	// c.svc = &services{
	// 	kvc:    eventstore.NewKVClient(conn),
	// 	mapc:   eventstore.NewMapClient(conn),
	// 	queuec: eventstore.NewQueueClient(conn),
	// }

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
