// Here is an example of how we can compose tasks
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

// This is a grouped task, it will allow us to build our own task
// leveraging other tasks
type getHostnameAndIPResult struct {
	SubResults []task.RemoteCommandResults
}

func (r getHostnameAndIPResult) String() string {
	return fmt.Sprintf("    hostname: %s    ip address: %s", r.SubResults[0].Stdout, r.SubResults[1].Stdout)
}

func (r *getHostnameAndIP) Run(ctx context.Context, host *gornir.Host) (interface{}, error) {
	// We call the first subtask and store the subresult
	res1, err := (&task.RemoteCommand{Command: "hostname"}).Run(ctx, host)
	if err != nil {
		return getHostnameAndIPResult{}, err
	}

	// We call the second subtask and store the subresult
	res2, err := (&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"}).Run(ctx, host)
	if err != nil {
		return getHostnameAndIPResult{}, err
	}
	return getHostnameAndIPResult{
		SubResults: []task.RemoteCommandResults{res1.(task.RemoteCommandResults), res2.(task.RemoteCommandResults)},
	}, nil
}

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

	rnr := runner.Sorted()

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	results, err := gr.RunSync(
		"Let's run a couple of commands",
		&getHostnameAndIP{},
	)
	if err != nil {
		log.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, true)
}
