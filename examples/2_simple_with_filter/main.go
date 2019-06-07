// Similar to the simple example but filtering the hosts
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
	// Instantiate a logger plugin
	log := logger.NewLogrus(false)

	// Load the inventory using the FromYAMLFile plugin
	file := "/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml"
	plugin := inventory.FromYAML{HostsFile: file}
	inv, err := plugin.Create()
	if err != nil {
		log.Fatal(err)
	}

	// define a function we will use to filter the hosts
	filter := func(ctx context.Context, h *gornir.Host) bool {
		return h.Hostname == "dev1.group_1" || h.Hostname == "dev4.group_2"
	}

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithFilter(filter)

	// Before calling Gornir.RunS we call Gornir.Filter and pass the function defined
	// above. This will narrow down the inventor to the hosts matching the filter
	results, err := gr.Filter(context.Background(), filter).RunSync(
		"What's my ip?",
		runner.Parallel(),
		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
	)
	if err != nil {
		log.Fatal(err)
	}

	output.RenderResults(os.Stdout, results, true)
}
