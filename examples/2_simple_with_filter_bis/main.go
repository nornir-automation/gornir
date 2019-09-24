// Similar to the simple_with_filter but leveraging included filters
package main

import (
	"context"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"
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

	// Before calling any method we created a filtered version of our gr object.
	// The original object is left unmodified. Now we can use this filtered
	// object to execute our tasks in a subset of our devices
	filteredGr := gr.Filter(filter)

	// Open an SSH connection towards the devices
	results, err := filteredGr.RunSync(
		context.Background(),
		&connection.SSHOpen{},
	)
	if err != nil {
		log.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, "Connecting to devices via ssh", true)

	// defer closing the SSH connection we just opened
	defer func() {
		results, err = filteredGr.RunSync(
			context.Background(),
			&connection.SSHClose{},
		)
		if err != nil {
			log.Fatal(err)
		}
		output.RenderResults(os.Stdout, results, "Close ssh connection", true)
	}()

	// Before calling Gornir.RunS we call Gornir.Filter and pass the function defined
	// above. This will narrow down the inventor to the hosts matching the filter
	results, err = filteredGr.Filter(filter).RunSync(
		context.Background(),
		&task.RemoteCommand{
			Command: "ip addr | grep \\/24 | awk '{ print $2 }'",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	output.RenderResults(os.Stdout, results, "What's my ip?", true)
}
