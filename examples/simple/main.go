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

	results, err := gr.RunS(
		"What's my hostname?",
		runner.Sequence(),
		&task.RemoteCommand{Command: "hostname"},
	)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(output.RenderResults(results))
}
