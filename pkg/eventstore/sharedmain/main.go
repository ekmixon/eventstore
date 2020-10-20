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

package sharedmain

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"

	"knative.dev/pkg/logging"
	"knative.dev/pkg/metrics"
	"knative.dev/pkg/signals"
)

// EventStoreServer interface
type EventStoreServer interface {
	Start(ctx context.Context) error
}

// StorageConstructor builds a event storage server that can be injected and run
type StorageConstructor func(ctx context.Context, env EnvConfigAccessor) EventStoreServer

// Main is the run wrapper for event storage servers
func Main(component string, envCtor EnvConfigConstructor, stCtor StorageConstructor) {
	MainWithContext(signals.NewContext(), component, envCtor, stCtor)
}

// MainWithContext is the run wrapper for event storage servers with context
func MainWithContext(ctx context.Context, component string, envCtor EnvConfigConstructor, stCtor StorageConstructor) {
	flag.Parse()

	env := envCtor()
	if err := envconfig.Process("", env); err != nil {
		log.Fatalf("Error processing env var: %s", err)
	}
	env.SetComponent(component)

	logger := env.GetLogger()
	defer flush(logger)
	ctx = logging.WithLogger(ctx, logger)

	msp := metrics.NewMemStatsAll()
	msp.Start(ctx, 30*time.Second)
	if err := view.Register(msp.DefaultViews()...); err != nil {
		logger.Fatal("Error exporting go memstats view: %v", zap.Error(err))
	}

	server := stCtor(ctx, env)

	logger.Info("Starting event storage server", zap.Any("server", server))

	if err := server.Start(ctx); err != nil {
		logger.Warn("start returned an error", zap.Error(err))
	}
}

func flush(logger *zap.SugaredLogger) {
	_ = logger.Sync()
	metrics.FlushExporter()
}
