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
	"fmt"
	"log"
	"net"
	"time"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"knative.dev/pkg/logging"

	"github.com/triggermesh/eventstore/pkg/eventstore/protob"
	"github.com/triggermesh/eventstore/pkg/eventstore/sharedmain"
)

const listenPort = "8080"

// ServerCtor creates an event store server
func ServerCtor(ctx context.Context, envAcc sharedmain.EnvConfigAccessor) sharedmain.EventStoreServer {
	// env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	return &inMemoryEventStore{
		logger: logger,
	}
}

type inMemoryEventStore struct {
	logger *zap.SugaredLogger
}

type storedValue struct {
	value   []byte
	expires time.Time
}

func (s *inMemoryEventStore) Start(ctx context.Context) error {
	s.logger.Info("Starting in memory event store")
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", listenPort))
	if err != nil {
		log.Fatal("failed to start listening: %s", err)
	}

	ims := &inMemoryServer{
		store: make(map[string]storedValue),
	}

	srv := grpc.NewServer()
	protob.RegisterEventStoreServer(srv, ims)
	srv.Serve(lis)

	return nil
}

// inMemoryServer is a ephemeral storage implementation of EventStore
type inMemoryServer struct {
	protob.UnimplementedEventStoreServer
	store map[string]storedValue
}

var _ protob.EventStoreServer = (*inMemoryServer)(nil)

func (s *inMemoryServer) Save(_ context.Context, sr *protob.SaveRequest) (*protob.SaveResponse, error) {
	if err := sr.Validate(); err != nil {
		return nil, err
	}

	t := tokenizer(sr.Location)

	s.store[t] = storedValue{
		value:   sr.GetValue(),
		expires: time.Now().Add(time.Duration(sr.GetTtl()) * time.Millisecond),
	}

	return &protob.SaveResponse{}, nil
}

func (s *inMemoryServer) Load(_ context.Context, lr *protob.LoadRequest) (*protob.LoadResponse, error) {
	if err := lr.Validate(); err != nil {
		return nil, err
	}

	t := tokenizer(lr.Location)
	v, ok := s.store[t]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "key %q not present at store", t)
	}

	if time.Now().After(v.expires) {
		delete(s.store, t)
		return nil, status.Errorf(codes.InvalidArgument, "key %q is expired", t)
	}

	return &protob.LoadResponse{Value: v.value}, nil
}

func (s *inMemoryServer) Delete(_ context.Context, dr *protob.DeleteRequest) (*protob.DeleteResponse, error) {
	if err := dr.Validate(); err != nil {
		return nil, err
	}

	t := tokenizer(dr.Location)
	_, ok := s.store[t]
	if ok {
		delete(s.store, t)
	}

	return &protob.DeleteResponse{}, nil
}

func tokenizer(location *protob.LocationType) string {
	var token string
	switch location.Scope.Type {
	case protob.ScopeChoice_Global:
		token = "global."
	case protob.ScopeChoice_Bridge:
		token = "bridge." + location.Scope.Bridge + "."
	case protob.ScopeChoice_Instance:
		token = "instance." + location.Scope.Bridge + "." + location.Scope.Instance + "."
	}

	return token + location.Key
}
