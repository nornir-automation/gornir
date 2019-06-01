package main

import (
	"fmt"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

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

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	rnr := runner.Sequence()
	err = gr.RunA(
		"What's my hostname?",
		rnr,
		&task.RemoteCommand{Command: "hostname"},
		results,
	)
	if err != nil {
		logger.Fatal(err)
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		_ = rnr.Wait()
		close(results)
	}()
	fmt.Println(output.RenderResults(results))
}
