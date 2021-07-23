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

type internalMap struct {
	*internalClient
}

type internalMapFields struct {
	*internalClient

	key string
}

var _ Map = (*internalMap)(nil)
var _ MapFields = (*internalMapFields)(nil)

// Set key/value at store
func (i *internalMap) New(ctx context.Context, key string, ttlSec int32) error {
	if i.svc.mapc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.NewMapRequest{
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

	_, err := i.svc.mapc.New(ctx, r)
	return err
}

// Get value from EventStore
func (i *internalMap) Fields(key string) MapFields {
	return &internalMapFields{
		internalClient: i.internalClient,
		key:            key,
	}
}

// Del Value from EventStore
func (i *internalMap) Del(ctx context.Context, key string) error {
	if i.svc.mapc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.DelMapRequest{
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

	_, err := i.svc.mapc.Del(ctx, r)
	return err
}

// Set map field.
func (i *internalMapFields) Set(ctx context.Context, key string, value []byte) error {
	if i.svc.mapc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.SetMapFieldRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
		Field: key,
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

	_, err := i.svc.mapc.FieldSet(ctx, r)
	return err
}

// Get map field.
func (i *internalMapFields) Get(ctx context.Context, key string) ([]byte, error) {
	if i.svc.mapc == nil {
		return nil, errors.New("EventStore client is not connected")
	}

	r := &eventstore.GetMapFieldRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
		Field: key,
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

	res, err := i.svc.mapc.FieldGet(ctx, r)
	if err != nil {
		return nil, err
	}

	return res.GetValue(), nil
}

// Del map field.
func (i *internalMapFields) Del(ctx context.Context, key string) error {
	if i.svc.mapc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.DelMapFieldRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
		Field: key,
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

	_, err := i.svc.mapc.FieldDel(ctx, r)
	return err
}

// Incr integer for field.
func (i *internalMapFields) Incr(ctx context.Context, key string, value int32) error {
	if i.svc.mapc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.IncrMapFieldRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
		Field: key,
		Incr:  value,
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

	_, err := i.svc.mapc.FieldIncr(ctx, r)
	return err
}

// Decr integer for field.
func (i *internalMapFields) Decr(ctx context.Context, key string, value int32) error {
	if i.svc.mapc == nil {
		return errors.New("EventStore client is not connected")
	}

	r := &eventstore.DecrMapFieldRequest{
		Location: &eventstore.LocationType{
			Scope: &eventstore.ScopeType{
				Bridge:   i.bridge,
				Instance: i.instance,
			},
			Key: i.key,
		},
		Field: key,
		Decr:  value,
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

	_, err := i.svc.mapc.FieldDecr(ctx, r)
	return err
}

// All elements in a map.
func (i *internalMapFields) All(ctx context.Context) (map[string][]byte, error) {
	if i.svc.mapc == nil {
		return nil, errors.New("EventStore client is not connected")
	}

	r := &eventstore.GetAllMapFieldsRequest{
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

	res, err := i.svc.mapc.GetFields(ctx, r)
	if err != nil {
		return nil, err
	}

	return res.GetValues(), nil
}

// Len for map.
func (i *internalMapFields) Len(ctx context.Context) (int, error) {
	if i.svc.mapc == nil {
		return 0, errors.New("EventStore client is not connected")
	}

	r := &eventstore.LenMapRequest{
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

	res, err := i.svc.mapc.Len(ctx, r)
	return int(res.GetLen()), err
}
