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

type internalKV struct {
	*internalClient
}

var _ KeyValue = (*internalKV)(nil)

// Set key/value at store
func (i *internalKV) Set(ctx context.Context, key string, value []byte, ttlSec int32) error {
	if i.svc.kvc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.SetKVRequest{
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
		r.Location.Scope.Type = eventstore.ScopeChoice_Instance
	case i.bridge != "":
		r.Location.Scope.Type = eventstore.ScopeChoice_Bridge
	default:
		r.Location.Scope.Type = eventstore.ScopeChoice_Global
	}

	if err := r.Validate(); err != nil {
		return err
	}

	_, err := i.svc.kvc.Set(ctx, r)
	return err
}

// Get value from EventStore
func (i *internalKV) Get(ctx context.Context, key string) ([]byte, error) {
	if i.svc.kvc == nil {
		return nil, errors.New("EventStore client is not connected")
	}

	r := &eventstore.GetKVRequest{
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
		return nil, err
	}

	res, err := i.svc.kvc.Get(ctx, r)
	if err != nil {
		return nil, err
	}

	return res.GetValue(), nil
}

// Del Value from EventStore
func (i *internalKV) Del(ctx context.Context, key string) error {
	if i.svc.kvc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.DelKVRequest{
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

	_, err := i.svc.kvc.Del(ctx, r)
	return err
}

// Incr integer for key.
func (i *internalKV) Incr(ctx context.Context, key string, incr int32) error {
	if i.svc.kvc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.IncrKVRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: key,
		},
		Incr: incr,
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

	_, err := i.svc.kvc.Incr(ctx, r)
	return err
}

// Decr integer for key.
func (i *internalKV) Decr(ctx context.Context, key string, decr int32) error {
	if i.svc.kvc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.DecrKVRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: key,
		},
		Decr: decr,
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

	_, err := i.svc.kvc.Decr(ctx, r)
	return err
}

// Locks key temporarily.
func (i *internalKV) Lock(ctx context.Context, key string, timeout int32) (string, error) {
	if i.svc.kvc == nil {
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

	res, err := i.svc.kvc.Lock(ctx, r)
	if err != nil {
		return "", err
	}

	return res.GetUnlock(), err
}

// Unlock key.
func (i *internalKV) Unlock(ctx context.Context, key string, unlock string) error {
	if i.svc.kvc == nil {
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

	_, err := i.svc.kvc.Unlock(ctx, r)
	return err
}
