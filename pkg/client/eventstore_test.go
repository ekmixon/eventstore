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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triggermesh/eventstore/pkg/protob"
	fakees "github.com/triggermesh/eventstore/pkg/protob/fake"
)

const (
	// operation requests
	tSave   = "Save"
	tLoad   = "Load"
	tDelete = "Delete"

	tBridge   = "test-bridge"
	tInstance = "test-instance"

	// request parameters

	tKey = "test-key"
)

var (
	tTTL        = int32(60)
	tValue      = []byte("test-value")
	tEmptyValue = []byte(nil)
)

func TestGlobalEventStoreClient(t *testing.T) {
	esc := newFakeClient()
	c := &client{esClient: esc}
	ctx := context.Background()

	expected := []fakees.Request{}

	client := c.Global()
	_ = client.SaveValue(ctx, tKey, []byte(tValue), tTTL)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Global, tSave, "", "", tKey, tValue, tTTL))
	_, _ = client.LoadValue(ctx, tKey)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Global, tLoad, "", "", tKey, tEmptyValue, 0))
	_ = client.DeleteValue(ctx, tKey)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Global, tDelete, "", "", tKey, tEmptyValue, 0))

	client = c.Bridge(tBridge)
	_ = client.SaveValue(ctx, tKey, []byte(tValue), tTTL)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Bridge, tSave, tBridge, "", tKey, tValue, tTTL))
	_, _ = client.LoadValue(ctx, tKey)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Bridge, tLoad, tBridge, "", tKey, tEmptyValue, 0))
	_ = client.DeleteValue(ctx, tKey)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Bridge, tDelete, tBridge, "", tKey, tEmptyValue, 0))

	client = c.Instance(tBridge, tInstance)
	_ = client.SaveValue(ctx, tKey, []byte(tValue), tTTL)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Instance, tSave, tBridge, tInstance, tKey, tValue, tTTL))
	_, _ = client.LoadValue(ctx, tKey)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Instance, tLoad, tBridge, tInstance, tKey, tEmptyValue, 0))
	_ = client.DeleteValue(ctx, tKey)
	expected = append(expected, expectedRequest(protob.ScopeChoice_Instance, tDelete, tBridge, tInstance, tKey, tEmptyValue, 0))

	requests := esc.GetRequests()
	assert.Equal(t, len(expected), len(requests), "Unexpected number of requests")

	for i := range requests {
		//nolint:govet
		assert.Equal(t, expected[i], requests[i],
			"unexpected request at %s/%s",
			expected[i].Location.Scope.Type.String(),
			expected[i].Operation)
	}
}

func expectedRequest(scope protob.ScopeChoice, operation, bridge, instance, key string, value []byte, ttl int32) fakees.Request {
	return fakees.Request{
		Operation: operation,
		Location: protob.LocationType{
			Scope: &protob.ScopeType{
				Type:     scope,
				Bridge:   bridge,
				Instance: instance,
			},
			Key: key,
		},
		TTL:   ttl,
		Value: value,
	}
}

func newFakeClient() fakees.EventStoreClient {

	esClient := fakees.NewEventStoreClientFake(
		fakees.WithSave(func(in *protob.SaveRequest) (*protob.SaveResponse, error) {
			if in.Location.Key == "return error" {
				return nil, errors.New("fake error")
			}
			return nil, nil
		}))

	return esClient
}
