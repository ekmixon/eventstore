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
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/triggermesh/eventstore/pkg/eventstore/protob"
)

var (
	serverAddr = flag.String("server_addr", "localhost:9090", "The server address in the format of host:port")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to dial %s: %v", *&serverAddr, err)
	}
	defer conn.Close()
	client := protob.NewEventStoreClient(conn)

	log.Println("Saving value")
	_, err = client.Save(context.Background(), &protob.SaveRequest{
		Location: &protob.LocationType{
			Scope: &protob.ScopeType{
				Type:     protob.ScopeChoice_Instance,
				Bridge:   "mybridge",
				Instance: "01234-deadbeef",
			},
		},
		Value: "saveme1",
		Ttl:   1000,
	})
	if err != nil {
		log.Fatalf("could not save: %v", err)
	}

	time.Sleep(2 * time.Second)

	log.Println("Loading value")
	lr, err := client.Load(context.Background(), &protob.LoadRequest{
		Location: &protob.LocationType{
			Scope: &protob.ScopeType{
				Type:     protob.ScopeChoice_Instance,
				Bridge:   "mybridge",
				Instance: "01234-deadbeef",
			},
		},
	})

	if err != nil {
		log.Fatalf("could not load: %v", err)
	}
	log.Printf("Loaded value: %s", lr.GetValue())

}
