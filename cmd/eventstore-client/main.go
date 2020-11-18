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

package main

import (
	"context"
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/triggermesh/eventstore/pkg/client"
)

var (
	server   = kingpin.Flag("server", "Event storage address.").Required().String()
	command  = kingpin.Flag("command", "One of load, save, delete.").Required().String()
	scope    = kingpin.Flag("scope", "One of Global, Bridge, Instance.").Required().String()
	bridge   = kingpin.Flag("bridge", "Bridge name, when scope is bridge or instance.").Default("").String()
	instance = kingpin.Flag("instance", "Instance ID, when scope is instance.").Default("").String()
	key      = kingpin.Flag("key", "Key for storing value.").Required().String()
	value    = kingpin.Flag("value", "Value to be stored.").Default("").String()
	ttl      = kingpin.Flag("ttl", "Stored value's time to live (seconds).").Default("5").Int32()

	timeout = kingpin.Flag("timeout", "Timeout for completing the operation.").Default("5s").Duration()
)

func main() {
	kingpin.Parse()

	// flag validation
	ctx := context.Background()

	c := client.New(*server, *timeout)
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("Failed to dial %s: %v", *server, err)
	}

	defer func() { _ = c.Disconnect() }()

	var sc client.Interface

	switch {
	case *scope == "Global":
		sc = c.Global()
	case *scope == "Bridge":
		sc = c.Bridge(*bridge)
	case *scope == "Instance":
		sc = c.Instance(*bridge, *instance)
	default:
		log.Fatalf("Incorrect scope %q", *scope)
	}

	// execution

	switch *command {
	case "load":
		r, err := sc.LoadValue(ctx, *key)
		if err != nil {
			log.Fatalf("could not load %q: %v", *key, err)
		}
		log.Printf("Loaded value (string): %s", string(r))

	case "save":
		err := sc.SaveValue(ctx, *key, []byte(*value), *ttl)
		if err != nil {
			log.Fatalf("could not save key %q: %v", *key, err)
		}
		log.Print("Saved.")

	case "delete":
		err := sc.DeleteValue(ctx, *key)
		if err != nil {
			log.Fatalf("could not delete key %q: %v", *key, err)
		}
		log.Print("Deleted.")

	default:
		kingpin.FatalUsage("Not valid command %q", *command)
	}
}
