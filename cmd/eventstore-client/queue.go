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
	"log"
	"time"

	"github.com/triggermesh/eventstore/pkg/client"
)

type QueueCmd struct {
	Set QueueNewCmd `cmd:"" help:"Create new queue at key"`
	Get QueueDelCmd `cmd:"" help:"Delete queue at key"`
}

type QueueNewCmd struct {
	Key string        `help:"Key where the value will be stored" required:""`
	TTL time.Duration `help:"Key's time to live (seconds)" default:"5s"`
}

type QueueDelCmd struct {
	Key string `help:"Queue key to delete" required:""`
}

func (kv *QueueNewCmd) Run(g *Globals) error {
	c := client.New(g.Server, g.Timeout)
	ctx := context.Background()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("Failed to dial %s: %v", g.Server, err)
	}

	defer func() { _ = c.Disconnect() }()

	return nil
}

func (kv *QueueDelCmd) Run(g *Globals) error {
	c := client.New(g.Server, g.Timeout)
	ctx := context.Background()
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("Failed to dial %s: %v", g.Server, err)
	}

	defer func() { _ = c.Disconnect() }()

	// TODO actual work
	return nil
}
