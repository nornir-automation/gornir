// this is the simplest example possible
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
	logger := logger.NewLogrus(false) // instantiate a logger plugin

	// Load the inventory using the FromYAMLFile plugin
	inventory, err := inventory.FromYAMLFile("/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	// Instantiate Gornir
	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	// Following call is going to execute the task over all the hosts using the runner.Parallel runner.
	// Said runner is going to handle the parallelization for us. Gornir.RunS is also going to block
	// until the runner has completed executing the task over all the hosts
	results, err := gr.RunS(
		"What's my ip?",
		runner.Parallel(),
		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
	)
	if err != nil {
		logger.Fatal(err)
	}

	// next call is going to print the result on screen
	output.RenderResults(os.Stdout, results, true)
}
