// Package gornir provides a pluggable framework with inventory management to help operate
// collections of devices.
// It's similar to https://github.com/nornir-automation/nornir/ but in Go.
//
// The goal is to be able to operate on many devices with little effort. For instance:
//
// 	package main
//
// 	import (
// 		"context"
// 		"os"
//
// 		"github.com/nornir-automation/gornir/pkg/gornir"
// 		"github.com/nornir-automation/gornir/pkg/plugins/connection"
// 		"github.com/nornir-automation/gornir/pkg/plugins/inventory"
// 		"github.com/nornir-automation/gornir/pkg/plugins/logger"
// 		"github.com/nornir-automation/gornir/pkg/plugins/output"
// 		"github.com/nornir-automation/gornir/pkg/plugins/runner"
// 		"github.com/nornir-automation/gornir/pkg/plugins/task"
// 	)
//
// 	func main() {
// 		log := logger.NewLogrus(false)
//
// 		file := "/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml"
// 		plugin := inventory.FromYAML{HostsFile: file}
// 		inv, err := plugin.Create()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(runner.Parallel())
//
// 		results, err := gr.RunSync(
// 			context.Background(),
// 			&connection.SSHOpen{},
// 		)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		// defer closing the SSH connection we just opened
// 		defer func() {
// 			results, err = gr.RunSync(
// 				context.Background(),
// 				&connection.SSHClose{},
// 			)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}()
//
// 		results, err = gr.RunSync(
// 			context.Background(),
// 			&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
// 		)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		output.RenderResults(os.Stdout, results, "What is my ip?", true)
// 	}
//
// would render:
//
//     # What is my ip?
//     @ dev5.no_group
//       - err: failed to retrieve connection: couldn't find connection
//
//     @ dev1.group_1
//       - stdout: 10.21.33.101/24
//
//       - stderr:
//     @ dev6.no_group
//       - stdout: 10.21.33.106/24
//
//       - stderr:
//     @ dev4.group_2
//       - stdout: 10.21.33.104/24
//
//       - stderr:
//     @ dev3.group_2
//       - stdout: 10.21.33.103/24
//
//       - stderr:
//     @ dev2.group_1
//       - stdout: 10.21.33.102/24
//
//       - stderr:
//
// You can see more examples here: https://github.com/nornir-automation/gornir/tree/master/examples
package gornir
