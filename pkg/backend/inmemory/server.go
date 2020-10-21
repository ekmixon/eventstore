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
	"sync"
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
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	defaultTTL := map[protob.ScopeChoice]int32{
		protob.ScopeChoice_Global:   env.DefaultGlobalTTL,
		protob.ScopeChoice_Bridge:   env.DefaultBridgeTTL,
		protob.ScopeChoice_Instance: env.DefaultInstanceTTL,
	}

	return &inMemoryEventStore{
		store:           make(map[string]storedValue),
		defaultTTL:      defaultTTL,
		expiredGCPeriod: env.DefaultExpiredGCPeriod,

		logger: logger,
	}
}

// inMemoryEventStore is an implementation of EventStore based on memory ephemeral backend.
type inMemoryEventStore struct {
	protob.UnimplementedEventStoreServer
	store           map[string]storedValue
	defaultTTL      map[protob.ScopeChoice]int32
	expiredGCPeriod int32

	mutex  sync.RWMutex
	logger *zap.SugaredLogger
}

type storedValue struct {
	value   []byte
	expires time.Time
}

// Start event store server
func (s *inMemoryEventStore) Start(ctx context.Context) error {
	s.logger.Info("Starting in memory event store")
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", listenPort))
	if err != nil {
		log.Fatal("failed to start listening: %s", err)
	}

	srv := grpc.NewServer()
	protob.RegisterEventStoreServer(srv, s)

	defer srv.GracefulStop()

	errCh := make(chan error)
	go func() {
		if err := srv.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Duration(s.expiredGCPeriod) * time.Second)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				s.deleteExpired()
			}
		}
	}()

	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}

var _ protob.EventStoreServer = (*inMemoryEventStore)(nil)

func (s *inMemoryEventStore) Save(_ context.Context, sr *protob.SaveRequest) (*protob.SaveResponse, error) {
	if err := sr.Validate(); err != nil {
		return nil, err
	}

	t := tokenizer(sr.Location)

	ttl := sr.GetTtl()
	if ttl == 0 {
		ttl = s.defaultTTL[sr.Location.Scope.Type]
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[t] = storedValue{
		value:   sr.GetValue(),
		expires: time.Now().Add(time.Duration(ttl) * time.Second),
	}

	return &protob.SaveResponse{}, nil
}

func (s *inMemoryEventStore) Load(_ context.Context, lr *protob.LoadRequest) (*protob.LoadResponse, error) {
	if err := lr.Validate(); err != nil {
		return nil, err
	}

	t := tokenizer(lr.Location)

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v, ok := s.store[t]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "key %q not present at store", t)
	}

	if time.Now().After(v.expires) {
		return nil, status.Errorf(codes.InvalidArgument, "key %q is expired", t)
	}

	return &protob.LoadResponse{Value: v.value}, nil
}

func (s *inMemoryEventStore) Delete(_ context.Context, dr *protob.DeleteRequest) (*protob.DeleteResponse, error) {
	if err := dr.Validate(); err != nil {
		return nil, err
	}

	t := tokenizer(dr.Location)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.store[t]
	if ok {
		delete(s.store, t)
	}

	return &protob.DeleteResponse{}, nil
}

func (s *inMemoryEventStore) deleteExpired() {

	// we dont mind if there are dirty reads, no need to lock
	expired := []string{}
	now := time.Now()
	for k, v := range s.store {
		if now.After(v.expires) {
			expired = append(expired, k)
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, k := range expired {
		delete(s.store, k)
	}
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
