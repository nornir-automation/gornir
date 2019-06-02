// Package provides a pluggable framework with inventory management to help operate collections of devices.
// It's similar to https://github.com/nornir-automation/nornir/ but in golang.
//
// The goal is to be able to operate on many devices with little effort. For instance:
//
//     package main
//
//     import (
//     	"os"
//
//     	"github.com/nornir-automation/gornir/pkg/gornir"
//     	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
//     	"github.com/nornir-automation/gornir/pkg/plugins/logger"
//     	"github.com/nornir-automation/gornir/pkg/plugins/output"
//     	"github.com/nornir-automation/gornir/pkg/plugins/runner"
//     	"github.com/nornir-automation/gornir/pkg/plugins/task"
//     )
//
//     func main() {
//     	logger := logger.NewLogrus(false)
//
//     	inventory, err := inventory.FromYAMLFile("/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml")
//     	if err != nil {
//     		logger.Fatal(err)
//     	}
//
//     	gr := &gornir.Gornir{
//     		Inventory: inventory,
//     		Logger:    logger,
//     	}
//
//     	results, err := gr.RunS(
//     		"What's my ip?",
//     		runner.Parallel(),
//     		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
//     	)
//     	if err != nil {
//     		logger.Fatal(err)
//     	}
//     	output.RenderResults(os.Stdout, results, true)
//     }
//
// would render:
//
//     # What's my ip?
//     @ dev5.no_group
//       - err: failed to dial: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain
//
//     @ dev1.group_1
//      * Stdout: 10.21.33.101/24
//
//      * Stderr:
//       - err: <nil>
//
//     @ dev2.group_1
//      * Stdout: 10.21.33.102/24
//
//      * Stderr:
//       - err: <nil>
//
//     @ dev3.group_2
//      * Stdout: 10.21.33.103/24
//
//      * Stderr:
//       - err: <nil>
//
//     @ dev4.group_2
//      * Stdout: 10.21.33.104/24
//
//      * Stderr:
//       - err: <nil>
//
// You can see more examples here: https://github.com/nornir-automation/gornir/tree/master/examples
package gornir
