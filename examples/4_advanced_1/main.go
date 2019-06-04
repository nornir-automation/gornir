// In this example we can see how we can call the runner asynchronously
package main

import (
	"context"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	// Instantiate a logger plugin.
	logger := logger.NewLogrus(false)
	// File where the inventory will be loaded from.
	file := "/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml"
	plugin := inventory.FromYAML{HostsFile: file}

	// Instantiate Gornir
	gr, err := gornir.Build(
		gornir.WithInventory(plugin),
		gornir.WithLogger(logger),
	)
	if err != nil {
		logger.Fatal(err)
	}

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	// We need to store the runner as we will need to check its completion later on
	// by calling rnr.Wait()
	rnr := runner.Parallel()

	// Gornir.RunAsync doesn't block so it's up to the user to check the runner is done
	err = gr.RunAsync(
		context.Background(),
		"What's my hostname?",
		rnr,
		&task.RemoteCommand{Command: "hostname"},
		results,
	)
	if err != nil {
		logger.Fatal(err)
	}

	// Next call will block until the runner is done
	rnr.Wait()

	close(results) // we need to close the channel or output.RenderResults will not finish
	output.RenderResults(os.Stdout, results, true)
}
