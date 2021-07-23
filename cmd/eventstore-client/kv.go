/*
Copyright (c) 2021 TriggerMesh Inc.

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

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/triggermesh/eventstore/pkg/client"
)

type KVCmd struct {
	Set  KVSetCmd  `cmd:"" help:"Set Key/Value"`
	Get  KVGetCmd  `cmd:"" help:"Get Value"`
	Del  KVDelCmd  `cmd:"" help:"Delete Key"`
	Incr KVIncrCmd `cmd:"" help:"Increase value"`
	Decr KVDecrCmd `cmd:"" help:"Increase value"`
}

type KVSetCmd struct {
	Value string        `help:"Value to be stored" required:""`
	TTL   time.Duration `help:"Key's time to live (seconds)" default:"5s"`
}

type KVGetCmd struct{}

type KVDelCmd struct{}

type KVIncrCmd struct {
	Incr int32 `help:"Value to be increased" default:"1"`
}

type KVDecrCmd struct {
	Decr int32 `help:"Value to be decreased" default:"1"`
}

func (kv *KVSetCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	if err := g.scopedClient(es).KV().Set(ctx, g.Key, []byte(kv.Value), int32(kv.TTL)); err != nil {
		return err
	}

	printDone()
	return nil
}

func (kv *KVGetCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	value, err := g.scopedClient(es).KV().Get(ctx, g.Key)
	if err != nil {
		return err
	}

	printKV(g.Key, string(value))
	return nil
}

func (kv *KVDelCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	if err := g.scopedClient(es).KV().Del(ctx, g.Key); err != nil {
		return err
	}

	printDone()
	return nil
}

func (kv *KVIncrCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	if err := g.scopedClient(es).KV().Incr(ctx, g.Key, kv.Incr); err != nil {
		return err
	}

	printDone()
	return nil
}

func (kv *KVDecrCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	if err := g.scopedClient(es).KV().Decr(ctx, g.Key, kv.Decr); err != nil {
		return err
	}

	printDone()
	return nil
}
