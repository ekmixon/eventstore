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

	eventstore "github.com/triggermesh/eventstore/pkg/protob"
)

type internalQueue struct {
	*internalClient
}

type internalQueueItems struct {
	*internalClient

	key string
}

var _ Queue = (*internalQueue)(nil)
var _ QueueItems = (*internalQueueItems)(nil)

// Set key/value at store
func (i *internalQueue) New(ctx context.Context, key string, ttlSec int32) error {
	if i.svc.queuec == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.NewQueueRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: key,
		},
		Ttl: ttlSec,
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return err
	}

	_, err := i.svc.queuec.New(ctx, r)
	return err
}

// Get value from EventStore
func (i *internalQueue) Items(key string) QueueItems {
	return &internalQueueItems{
		internalClient: i.internalClient,
		key:            key,
	}
}

// Del Value from EventStore
func (i *internalQueue) Del(ctx context.Context, key string) error {
	if i.svc.queuec == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.DelQueueRequest{
		Location: &eventstore.LocationType{
			Key: key,
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			}}}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return err
	}

	_, err := i.svc.queuec.Del(ctx, r)
	return err
}

// Locks key temporarily.
func (i *internalQueue) Lock(ctx context.Context, key string, timeout int32) (string, error) {
	if i.svc.queuec == nil {
		return "", errors.New("EventStore client is not connected")
	}

	r := &eventstore.LockRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: key,
		},
		Timeout: timeout,
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return "", err
	}

	res, err := i.svc.queuec.Lock(ctx, r)
	if err != nil {
		return "", err
	}
	return res.Unlock, err
}

// Unlock key.
func (i *internalQueue) Unlock(ctx context.Context, key string, unlock string) error {
	if i.svc.queuec == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.UnlockRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: key,
		},
		Unlock: unlock,
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return err
	}

	_, err := i.svc.queuec.Unlock(ctx, r)
	return err
}

// Push item to the queue.
func (i *internalQueueItems) Push(ctx context.Context, value []byte) error {
	if i.svc.queuec == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.PushQueueRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
		Value: value,
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return err
	}

	_, err := i.svc.queuec.Push(ctx, r)
	return err
}

// Pop item, removing it from the queue
func (i *internalQueueItems) Pop(ctx context.Context) ([]byte, error) {
	if i.svc.queuec == nil {
		return nil, errors.New("EventStore client is not connected")
	}

	r := &eventstore.PopQueueRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	res, err := i.svc.queuec.Pop(ctx, r)
	if err != nil {
		return nil, err
	}

	return res.GetValue(), nil
}

// Peek item, keep it in the queue
func (i *internalQueueItems) Peek(ctx context.Context) ([]byte, error) {
	if i.svc.queuec == nil {
		return nil, errors.New("EventStore client is not connected")
	}

	r := &eventstore.PeekQueueRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	res, err := i.svc.queuec.Peek(ctx, r)
	return res.GetValue(), err
}

// Index item, removing it from the queue.
func (i *internalQueueItems) Index(ctx context.Context, index int32) ([]byte, error) {
	if i.svc.queuec == nil {
		return nil, errors.New("EventStore client is not connected")
	}

	r := &eventstore.IndexQueueRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
		Index: index,
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	res, err := i.svc.queuec.Index(ctx, r)
	if err != nil {
		return nil, err
	}

	return res.GetValue(), nil
}

// All elements in a map.
func (i *internalQueueItems) All(ctx context.Context) ([][]byte, error) {
	if i.svc.queuec == nil {
		return nil, errors.New("EventStore client is not connected")
	}

	r := &eventstore.GetAllQueuesRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	res, err := i.svc.queuec.GetAll(ctx, r)
	if err != nil {
		return nil, err
	}

	return res.GetValues(), nil
}

// Len for map.
func (i *internalQueueItems) Len(ctx context.Context) (int, error) {
	if i.svc.queuec == nil {
		return 0, errors.New("EventStore client is not connected")
	}

	r := &eventstore.LenQueueRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
	}

	switch {
	case i.instance != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return 0, err
	}

	res, err := i.svc.queuec.Len(ctx, r)
	return int(res.GetLen()), err
}
