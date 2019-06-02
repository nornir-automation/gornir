// In this example we can see how we can call the runner asynchronously
package main

import (
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	logger := logger.NewLogrus(false)

	inventory, err := inventory.FromYAMLFile("/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	// We need to store the runner as we will need to check its completion later on
	// by calling rnr.Wait()
	rnr := runner.Parallel()

	// Gornir.RunA doesn't block so it's up to the user to check the runner is done
	err = gr.RunA(
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
