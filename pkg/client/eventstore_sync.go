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

type internalSync struct {
	*internalClient
}

var _ Sync = (*internalSync)(nil)

// Locks key temporarily.
func (i *internalSync) Lock(ctx context.Context, key string, timeout int32) (string, error) {
	if i.svc.syncc == nil {
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

	res, err := i.svc.syncc.Lock(ctx, r)
	if err != nil {
		return "", err
	}

	return res.GetUnlock(), err
}

// Unlock key.
func (i *internalSync) Unlock(ctx context.Context, key string, unlock string) error {
	if i.svc.syncc == nil {
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

	_, err := i.svc.syncc.Unlock(ctx, r)
	return err
}
