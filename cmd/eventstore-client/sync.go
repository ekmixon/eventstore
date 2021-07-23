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

type SyncCmd struct {
	Lock   LockCmd   `cmd:"" help:"Lock key for exclusive access"`
	Unlock UnlockCmd `cmd:"" help:"Unlock key"`
}

type LockCmd struct {
	UnlockTimeout int32 `help:"Timeout before automatically unlocking (seconds)" required:""`
}

type UnlockCmd struct {
	Unlock string `help:"Unlock string" required:""`
}

func (s *LockCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	unlock, err := g.scopedClient(es).Sync().Lock(ctx, g.Key, s.UnlockTimeout)
	if err != nil {
		return err
	}

	printKV("unlock", unlock)
	return nil
}

func (s *UnlockCmd) Run(g *Globals) error {
	es := client.New(g.Server, g.Timeout)
	ctx := context.Background()

	if err := es.Connect(ctx); err != nil {
		return fmt.Errorf("failed to dial %s: %v", g.Server, err)
	}
	defer func() { _ = es.Disconnect() }()

	if err := g.scopedClient(es).Sync().Unlock(ctx, g.Key, s.Unlock); err != nil {
		return err
	}

	printDone()
	return nil
}
