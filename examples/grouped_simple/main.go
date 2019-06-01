package main

import (
	"fmt"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

type checkMemoryAndCPU struct {
}

func (c checkMemoryAndCPU) Run(ctx gornir.Context, wg *sync.WaitGroup, jobResult chan *gornir.JobResult) {
	result := gornir.NewJobResult(ctx)
	defer wg.Done()

	sr := make(chan *gornir.JobResult, 1)
	swg := &sync.WaitGroup{}
	swg.Add(2)
	(&task.RemoteCommand{Command: "free -m"}).Run(ctx, swg, sr)
	result.AddSubResult(<-sr)
	(&task.RemoteCommand{Command: "uptime"}).Run(ctx, swg, sr)
	result.AddSubResult(<-sr)

	jobResult <- result
}

func main() {
	logger := logger.NewLogrus()

	inventory, err := inventory.FromYAMLFile("hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	results, err := gr.RunS(
		"Let's run a couple of commands",
		runner.Sequence(),
		&checkMemoryAndCPU{},
	)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(output.RenderResults(results))
}
