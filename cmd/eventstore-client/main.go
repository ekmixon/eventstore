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
	"time"

	"github.com/alecthomas/kong"
)

type Globals struct {
	Server   string `help:"Event storage address" required:""`
	Scope    string `help:"Storage scope" enum:"global,bridge,instance"`
	Bridge   string `help:"Bridge name, when scope is bridge or instance"`
	Instance string `help:"Instance ID, when scope is instance"`

	Timeout time.Duration `help:"Timeout for completing the operation" default:"5s"`
}

type Cli struct {
	Globals

	Kv KVCmd `cmd help:"KV EventStore"`

	// Kv struct {
	// 	Command string `help:"KV storage command."`
	// 	// Force     bool   `help:"Force removal." short:"f"`
	// 	// Recursive bool   `help:"Recursively remove files." short:"r"`

	// 	// Paths []string `arg:"" help:"Paths to remove." type:"path" name:"path"`
	// } `cmd:"" help:"EventStore KV."`

	Map struct {
		Paths []string `arg:"" optional:"" help:"Paths to list." type:"path"`
	} `cmd:"" help:"List paths."`

	Queue struct {
		Paths []string `arg:"" optional:"" help:"Paths to list." type:"path"`
	} `cmd:"" help:"List paths."`
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
