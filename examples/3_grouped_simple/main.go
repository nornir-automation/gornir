// Here is an example of how we can compose tasks
package main

import (
	"context"
	"os"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

// This is a grouped task, it will allow us to build our own task
// leveraging other tasks
type checkMemoryAndCPU struct {
}

func (c *checkMemoryAndCPU) Run(ctx context.Context, wg *sync.WaitGroup, taskParameters *gornir.TaskParameters, jobResult chan *gornir.JobResult) {
	// We instantiate a new object
	result := gornir.NewJobResult(ctx, taskParameters)

	defer wg.Done() // flag as completed

	// channel to store the subresults
	sr := make(chan *gornir.JobResult, 1)

	// We are going to execute two tasks so we need a sync.WaitGroup with two tokens
	swg := &sync.WaitGroup{}
	swg.Add(2)

	// We call the first subtask and store the subresult
	(&task.RemoteCommand{Command: "free -m"}).Run(ctx, swg, taskParameters, sr)
	result.AddSubResult(<-sr)

	// We call the second subtask and store the subresult
	(&task.RemoteCommand{Command: "uptime"}).Run(ctx, swg, taskParameters, sr)
	result.AddSubResult(<-sr)

	jobResult <- result
}

func main() {
	// main is business as usual
	logger := logger.NewLogrus(false)

	inventory, err := inventory.FromYAMLFile("/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	results, err := gr.RunS(
		"Let's run a couple of commands",
		runner.Parallel(),
		&checkMemoryAndCPU{},
	)
	if err != nil {
		logger.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, true)
}
