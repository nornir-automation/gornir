// Similar to the simple_with_filter but leveraging included filters
package main

import (
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	f "github.com/nornir-automation/gornir/pkg/plugins/filter"
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

	// this time our filter is composed from various FilterFunc
	filter := f.Or(f.WithHostname("dev1.group_1"), f.WithHostname("dev4.group_2"))

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(runner.Sorted())

	// Before calling Gornir.RunS we call Gornir.Filter and pass the function defined
	// above. This will narrow down the inventor to the hosts matching the filter
	results, err := gr.Filter(filter).RunSync(
		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
	)
	if err != nil {
		log.Fatal(err)
	}

	output.RenderResults(os.Stdout, results, "What's my ip?", true)
}
