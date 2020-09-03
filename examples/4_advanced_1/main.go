// In this example we can see how we can call the runner asynchronously
package main

import (
	"context"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	// Instantiate a logger plugin
	log := logger.NewLogrus(false)

	// Load the inventory using the FromYAMLFile plugin
	file := "/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml"
	plugin := inventory.FromYAML{HostsFile: file}
	inv, err := plugin.Create()
	if err != nil {
		log.Fatal(err)
	}

	// We need to store the runner as we will need to check its completion later on
	// by calling rnr.Wait()
	rnr := runner.Sorted()

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	// Open an SSH connection towards the devices
	results, err := gr.RunSync(
		context.Background(),
		&connection.SSHOpen{},
	)
	if err != nil {
		log.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, "Connecting to devices via ssh", true)

	// defer closing the SSH connection we just opened
	defer func() {
		results, err = gr.RunSync(
			context.Background(),
			&connection.SSHClose{},
		)
		if err != nil {
			log.Fatal(err)
		}
		output.RenderResults(os.Stdout, results, "Close ssh connection", true)
	}()

	// We need a channel to store the results
	res := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	// Gornir.RunAsync doesn't block so it's up to the user to check the runner is done
	err = gr.RunAsync(
		context.Background(),
		&task.RemoteCommand{Command: "hostname"},
		res,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Next call will block until the runner is done
	rnr.Wait()

	close(res) // we need to close the channel or output.RenderResults will not finish
	output.RenderResults(os.Stdout, res, "What's my hostname?", true)
}
