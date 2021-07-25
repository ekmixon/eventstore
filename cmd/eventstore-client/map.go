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

type MapCmd struct {
	New MapNewCmd `cmd:"" help:"Create new map at key"`
	Del MapDelCmd `cmd:"" help:"Delete map at key"`

	FieldSet  MapFieldSetCmd  `cmd:"" help:"Set value at field"`
	FieldGet  MapFieldGetCmd  `cmd:"" help:"Get value from field"`
	FieldDel  MapFieldDelCmd  `cmd:"" help:"Delete value from field"`
	FieldIncr MapFieldIncrCmd `cmd:"" help:"Increase value at field"`
	FieldDecr MapFieldDecrCmd `cmd:"" help:"Decrease value from field"`

	AllItems MapAllItemsCmd `cmd:"" help:"Get all items at map"`
	Len      MapLenCmd      `cmd:"" help:"Get map length"`
}

type MapNewCmd struct {
	TTL int32 `help:"Key's time to live (seconds)" default:"5"`
}

type MapDelCmd struct{}

type MapFieldSetCmd struct {
	Field string `help:"Field at map" required:""`
	Value string `help:"Value to be set at field" required:""`
}

type MapFieldGetCmd struct {
	Field string `help:"Field at map" required:""`
}

type MapFieldDelCmd struct {
	Field string `help:"Field at map" required:""`
}

type MapFieldIncrCmd struct {
	Field string `help:"Field at map" required:""`
	Incr  int32  `help:"Value to be increased" default:"1"`
}

type MapFieldDecrCmd struct {
	Field string `help:"Field at map" required:""`
	Decr  int32  `help:"Value to be decreased" default:"1"`
}

type MapAllItemsCmd struct{}
type MapLenCmd struct{}

func (s *MapNewCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Map().New(ctx, g.Key, s.TTL)
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (kv *MapDelCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Map().Del(ctx, g.Key)
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *MapFieldSetCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Map().Fields(g.Key).Set(ctx, s.Field, []byte(s.Value))
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *MapFieldGetCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Map().Fields(g.Key).Get(ctx, s.Field)
	if err != nil {
		return err
	}

	printKV(fmt.Sprintf("%s[%s]", g.Key, s.Field), res)
	return nil
}

func (s *MapFieldDelCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Map().Fields(g.Key).Del(ctx, s.Field)
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *MapFieldIncrCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Map().Fields(g.Key).Incr(ctx, s.Field, s.Incr)
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *MapFieldDecrCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	err := g.scopedClient(es).Map().Fields(g.Key).Decr(ctx, s.Field, s.Decr)
	if err != nil {
		return err
	}

	printDone()
	return nil
}

func (s *MapAllItemsCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Map().Fields(g.Key).All(ctx)
	if err != nil {
		return err
	}

	printMap("items", res)
	return nil
}

func (s *MapLenCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	res, err := g.scopedClient(es).Map().Fields(g.Key).Len(ctx)
	if err != nil {
		return err
	}

	printKV("len", res)
	return nil
}
