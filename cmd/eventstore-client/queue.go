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

	"github.com/triggermesh/eventstore/pkg/client"
)

type QueueCmd struct {
	New QueueNewCmd `cmd:"" help:"Create new queue"`
	Del QueueDelCmd `cmd:"" help:"Delete queue at key"`

	Push  QueuePushCmd  `cmd:"" help:"Push item to queue"`
	Index QueueIndexCmd `cmd:"" help:"Retrieve element at index"`
	Pop   QueuePopCmd   `cmd:"" help:"Extract and remove element at head"`
	Peek  QueuePeekCmd  `cmd:"" help:"Extract element at head wihtout removing"`

	AllItems QueueAllItemsCmd `cmd:"" help:"Get all items at queue"`
	Len      QueueLenCmd      `cmd:"" help:"Get queue length"`
}

type QueueNewCmd struct {
	TTL int32 `help:"Key's time to live (seconds)" default:"5"`
}

type QueueDelCmd struct {
}

type QueuePushCmd struct {
	Value string `help:"Value to be pushed" required:""`
}

type QueueIndexCmd struct {
	Index int32 `help:"Index for the queue element to retrieve" required:""`
}

type QueuePopCmd struct{}
type QueuePeekCmd struct{}
type QueueAllItemsCmd struct{}
type QueueLenCmd struct{}

func (s *QueueNewCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Queue().New(ctx, g.Key, s.TTL)
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *QueueDelCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Queue().Del(ctx, g.Key)
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *QueuePushCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Queue().Items(g.Key).Push(ctx, []byte(s.Value))
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *QueueIndexCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Queue().Items(g.Key).Index(ctx, s.Index)
	if err != nil {
		return err
	}

	printKV(fmt.Sprintf("%s[%d]", g.Key, s.Index), res)
	return nil
}

func (s *QueuePopCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Queue().Items(g.Key).Pop(ctx)
	if err != nil {
		return err
	}

	printKV(g.Key, res)
	return nil
}

func (s *QueuePeekCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Queue().Items(g.Key).Peek(ctx)
	if err != nil {
		return err
	}

	printKV(g.Key, res)
	return nil
}

func (s *QueueAllItemsCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Queue().Items(g.Key).All(ctx)
	if err != nil {
		return err
	}

	printList("items", res)
	return nil
}

func (s *QueueLenCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Queue().Items(g.Key).Len(ctx)
	if err != nil {
		return err
	}

	printKV("len", res)
	return nil
}
