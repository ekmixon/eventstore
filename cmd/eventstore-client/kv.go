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
	Set KVSetCmd `cmd:"" help:"Set Key/Value"`
	Get KVGetCmd `cmd:"" help:"Set Key/Value"`
	// DetachKeys string `help:"Override the key sequence for detaching a container"`
	// NoStdin    bool   `help:"Do not attach STDIN"`
	// SigProxy   bool   `help:"Proxy all received signals to the process" default:"true"`

	// Container string `arg required help:"Container ID to attach to."`
}

type KVSetCmd struct {
	Key   string        `help:"Key where the value will be stored" required:""`
	Value string        `help:"Value to be stored" required:""`
	TTL   time.Duration `help:"Key's time to live (seconds)" default:"5s"`
}

// type KVSet struct {
// 	key   = kingpin.Flag("key", "Key for storing value.").Required().String()
// 	value = kingpin.Flag("value", "Value to be stored.").Default("").String()
// 	ttl   = kingpin.Flag("ttl", "Stored value's time to live (seconds).").Default("5").Int32()
// }

type KVGetCmd struct {
}

func (kv *KVSetCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	c := g.scopedClient(es).KV()
	return c.Set(ctx, kv.Key, []byte(kv.Value), int32(kv.TTL))
}

func (kv *KVGetCmd) Run(g *Globals) error {
	fmt.Printf("KVGet Scope: %s\n", g.Scope)
	return nil
}
