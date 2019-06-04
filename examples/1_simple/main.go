// this is the simplest example possible
package main

import (
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
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
	
	// Instantiate Gornir
	gr, err := gornir.Build(
		gornir.WithInventory(file),
		gornir.WithLogger(logger),
	)	
	if err != nil {
		logger.Fatal(err)
	}

	// Following call is going to execute the task over all the hosts using the runner.Parallel runner.
	// Said runner is going to handle the parallelization for us. Gornir.RunS is also going to block
	// until the runner has completed executing the task over all the hosts
	results, err := gr.RunSync(
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
