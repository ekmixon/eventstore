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

	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/triggermesh/eventstore/pkg/eventstore/protob"
)

var (
	server   = kingpin.Flag("server", "Event storage address.").Required().String()
	command  = kingpin.Flag("command", "One of load, save, delete.").Required().String()
	scope    = kingpin.Flag("scope", "One of Global, Bridge, Instance.").Required().String()
	bridge   = kingpin.Flag("bridge", "Bridge name, when scope is bridge or instance.").String()
	instance = kingpin.Flag("instance", "Instance ID, when scope is instance.").String()
	key      = kingpin.Flag("key", "Key for storing value.").Required().String()
	value    = kingpin.Flag("value", "Value to be stored.").String()
	ttl      = kingpin.Flag("ttl", "Stored value's time to live (ms).").Int32()
)

func main() {
	kingpin.Parse()

	// validation

	if *command != "load" && *command != "save" && *command != "delete" {
		kingpin.FatalUsage("not valid command %q", *command)
	}

	sc, ok := protob.ScopeChoice_value[*scope]
	if !ok {
		kingpin.FatalUsage("not valid scope %q", *scope)
	}
	scopeTypeChoice := protob.ScopeChoice(sc)

	if *scope != "global" && (bridge == nil || *bridge == "") {
		kingpin.FatalUsage("bridge and instance scope types needs to inform bridge %q")
	}

	if *scope == "instance" && (instance == nil || *instance == "") {
		kingpin.FatalUsage("instance scope types needs to inform instance")
	}

	location := &protob.LocationType{
		Scope: &protob.ScopeType{
			Type:     scopeTypeChoice,
			Bridge:   *bridge,
			Instance: *instance,
		},
	}

	// connection

	conn, err := grpc.Dial(*server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to dial %s: %v", *server, err)
	}

	defer conn.Close()
	client := protob.NewEventStoreClient(conn)
	ctx := context.Background()

	// execution

	switch *command {
	case "load":
		lr, err := client.Load(ctx, &protob.LoadRequest{Location: location})
		if err != nil {
			log.Fatalf("could not load: %v", err)
		}
		log.Printf("Loaded value: %s", lr.GetValue())

	case "save":
		if value == nil {
			kingpin.FatalUsage("value not informed")
		}

		if ttl == nil {
			kingpin.FatalUsage("TTL not informed")
		}

		_, err = client.Save(ctx, &protob.SaveRequest{
			Location: location,
			Value:    []byte(*value),
			Ttl:      *ttl,
		})
		if err != nil {
			log.Fatalf("could not save: %v", err)
		}
		log.Print("Saved.")

	case "delete":
		_, err = client.Delete(ctx, &protob.DeleteRequest{Location: location})
		if err != nil {
			log.Fatalf("could not delete value: %v", err)
		}
		log.Println("Deleted.")

	default:
		kingpin.FatalUsage("not valid command %q", *command)
	}
}
