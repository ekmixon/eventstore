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

package inmemory

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	tlog "knative.dev/pkg/logging/testing"

	"github.com/triggermesh/eventstore/pkg/eventstore/protob"
)

const (
	tDefaultGlobalTTL   = 10
	tDefaultBridgeTTL   = 5
	tDefaultInstanceTTL = 2
	tExpiredGCPeriod    = 5
)

type tOperationType int
type tLocationType int

const (
	loadOp = iota
	saveOp
	deleteOp
)

type tOperation struct {
	// operation instructions
	opType   tOperationType
	location protob.LocationType
	ttl      int32
	value    []byte

	// wait after operation (for TTL tests)
	sleepSeconds int

	// expected value
	expecErr string
}

func TestInMemoryStore(t *testing.T) {

	// Cases
	testCases := map[string]struct {
		operations []tOperation
	}{
		"global operations": {
			operations: []tOperation{
				{
					// store key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),
					value:    []byte("val1"),
				},
				{
					// retrieve key1
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),

					value: []byte("val1"),
				},
				{
					// overwrite key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),
					value:    []byte("val2"),
				},
				{
					// retrieve key1
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),

					value: []byte("val2"),
				},
				{
					// store key2 TTL 1 sec
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Global, "key2"),

					value: []byte("val3"),
					ttl:   1,
				},
				{
					// retrieve key2
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Global, "key2"),

					value:        []byte("val3"),
					sleepSeconds: 1,
				},
				{
					// delete key1
					opType:   deleteOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),
				},
				{
					// retrieve key1
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),

					expecErr: `rpc error: code = InvalidArgument desc = key "global.key1" not present at store`,
				},
				{
					// retrieve key2
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Global, "key2"),

					expecErr: `rpc error: code = InvalidArgument desc = key "global.key2" is expired`,
				},
			},
		},

		"multibucket operations": {
			operations: []tOperation{
				{
					// store global key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),
					value:    []byte("val1"),
				},
				{
					// store bridge key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Bridge, "key1", withBridge("bridgeA")),
					value:    []byte("val2"),
				},
				{
					// store instance key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Instance, "key1", withBridge("bridgeA"), withInstance("instanceA")),
					value:    []byte("val3"),
				},
				{
					// retrieve global key1
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1"),

					value: []byte("val1"),
				},
				{
					// retrieve bridge key1
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Bridge, "key1", withBridge("bridgeA")),

					value: []byte("val2"),
				},
				{
					// retrieve instance key1
					opType:   loadOp,
					location: *newLocation(protob.ScopeChoice_Instance, "key1", withBridge("bridgeA"), withInstance("instanceA")),

					value: []byte("val3"),
				},
			},
		},

		"error global with bridge info": {
			operations: []tOperation{
				{
					// store bridge key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1", withBridge("bridgeA")),
					value:    []byte("val2"),

					expecErr: "rpc error: code = Unknown desc = global scope should not inform bridge nor instance",
				},
			},
		},

		"error global with instance info": {
			operations: []tOperation{
				{
					// store bridge key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Global, "key1", withInstance("instanceA")),
					value:    []byte("val2"),

					expecErr: "rpc error: code = Unknown desc = global scope should not inform bridge nor instance",
				},
			},
		},

		"error bridge without bridge info": {
			operations: []tOperation{
				{
					// store bridge key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Bridge, "key1"),
					value:    []byte("val2"),

					expecErr: "rpc error: code = Unknown desc = bridge scope needs the bridge identifier to be informed",
				},
			},
		},

		"error instance without bridge info": {
			operations: []tOperation{
				{
					// store bridge key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Instance, "key1", withInstance("instanceA")),
					value:    []byte("val2"),

					expecErr: "rpc error: code = Unknown desc = instance scope needs bridge and instance identifiers to be informed",
				},
			},
		},

		"error instance without instance info": {
			operations: []tOperation{
				{
					// store bridge key1
					opType:   saveOp,
					location: *newLocation(protob.ScopeChoice_Instance, "key1", withBridge("bridgeA")),
					value:    []byte("val2"),

					expecErr: "rpc error: code = Unknown desc = instance scope needs bridge and instance identifiers to be informed",
				},
			},
		},
	}

	ttl := map[protob.ScopeChoice]int32{
		protob.ScopeChoice_Global:   tDefaultGlobalTTL,
		protob.ScopeChoice_Bridge:   tDefaultBridgeTTL,
		protob.ScopeChoice_Instance: tDefaultInstanceTTL,
	}

	im := &inMemoryEventStore{
		store:           make(map[string]storedValue),
		defaultTTL:      ttl,
		expiredGCPeriod: tExpiredGCPeriod,

		logger: tlog.TestLogger(t),
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(im)))
	if err != nil {
		t.Fatalf("Error starting gRCP buffered server")
	}

	defer conn.Close()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			// Reset store. Cases should be run sequentially.
			im.store = make(map[string]storedValue)

			for i := range tc.operations {
				op := &tc.operations[i]
				client := protob.NewEventStoreClient(conn)

				switch op.opType {
				case saveOp:
					t.Logf("Saving %s at %+v/%s", op.value, op.location.Scope, op.location.Key)
					_, err := client.Save(ctx, &protob.SaveRequest{
						Location: &op.location,
						Ttl:      op.ttl,
						Value:    op.value,
					})

					assert.Equal(t, op.expecErr, errorToString(err))

				case loadOp:
					t.Logf("Loading from %+v/%s", op.location.Scope, op.location.Key)
					lr, err := client.Load(ctx, &protob.LoadRequest{
						Location: &op.location,
					})

					require.Equal(t, op.expecErr, errorToString(err))

					var value []byte
					if lr != nil {
						value = lr.Value
					}
					assert.Equal(t, value, op.value)

				case deleteOp:
					t.Logf("Deleting key %+v/%s", op.location.Scope, op.location.Key)
					_, err := client.Delete(ctx, &protob.DeleteRequest{
						Location: &op.location,
					})

					assert.Equal(t, op.expecErr, errorToString(err))
				}

				if op.sleepSeconds != 0 {
					time.Sleep(time.Duration(op.sleepSeconds) * time.Second)
				}
			}
		})
	}
}

func errorToString(e error) string {
	if e == nil {
		return ""
	}
	return e.Error()
}

type grpcDialer func(context.Context, string) (net.Conn, error)

func dialer(im *inMemoryEventStore) grpcDialer {
	buff := bufconn.Listen(1024 * 1024)

	srv := grpc.NewServer()
	protob.RegisterEventStoreServer(srv, im)

	go func() {
		if err := srv.Serve(buff); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return buff.Dial()
	}
}

type locationOpts func(*protob.LocationType)

func newLocation(scope protob.ScopeChoice, key string, opts ...locationOpts) *protob.LocationType {
	l := &protob.LocationType{
		Scope: &protob.ScopeType{
			Type: scope,
		},
		Key: key,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func withBridge(b string) locationOpts {
	return func(l *protob.LocationType) {
		l.Scope.Bridge = b
	}
}

func withInstance(i string) locationOpts {
	return func(l *protob.LocationType) {
		l.Scope.Instance = i
	}
}
