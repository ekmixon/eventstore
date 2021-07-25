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
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alecthomas/kong"

	"github.com/triggermesh/eventstore/pkg/client"
)

type Globals struct {
	Server   string `help:"Event storage address" required:""`
	Scope    string `help:"Storage scope" enum:"global,bridge,instance"`
	Bridge   string `help:"Bridge name, when scope is bridge or instance"`
	Instance string `help:"Instance ID, when scope is instance"`
	Key      string `help:"Storage Key" required:""`

	Timeout time.Duration `help:"Timeout for completing the operation" default:"5s"`
}

type Cli struct {
	Globals

	Kv    KVCmd    `cmd:"" help:"KV store"`
	Queue QueueCmd `cmd:"" help:"Queue store"`
	Map   MapCmd   `cmd:"" help:"Map store"`
	Sync  SyncCmd  `cmd:"" help:"Lock and unlock keys"`
}

func main() {
	cli := Cli{}
	ctx := kong.Parse(&cli,
		kong.Name("eventstore-client"),
		kong.Description("EventStore command line utility."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}

func (c *Cli) Validate() error {
	switch c.Scope {
	case "global":
		if c.Bridge != "" || c.Instance != "" {
			return errors.New("global scope does not need bridge or instance informed")
		}

	case "bridge":
		if c.Bridge == "" {
			return errors.New("bridge scope need the bridge identifier informed")
		}
		if c.Instance != "" {
			return errors.New("bridge scope does not need instance informed")
		}

	case "instance":
		if c.Bridge == "" || c.Instance == "" {
			return errors.New("instance scope needs bridge and instance identifiers informed")
		}

	default:
		return fmt.Errorf("unknown scope %q", c.Scope)
	}

	return nil
}

func (g *Globals) scopedClient(c client.EventStore) client.Interface {
	switch g.Scope {
	case "global":
		return c.Global()

	case "bridge":
		return c.Bridge(g.Bridge)

	case "instance":
		return c.Instance(g.Bridge, g.Instance)
	}
	return nil
}

func printKV(key string, value interface{}) {
	log.Printf("%s: %s\n", key, value)
}

func printList(key string, values [][]byte) {
	log.Printf("%s:\n", key)
	for i := range values {
		log.Printf("\t%d: %s\n", i, string(values[i]))
	}
}

func printMap(key string, values map[string][]byte) {
	log.Printf("%s:\n", key)
	for k, v := range values {
		log.Printf("\t%s: %s\n", k, string(v))
	}
}

func printDone() {
	log.Println("done")
}
