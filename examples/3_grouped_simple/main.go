// Here is an example of how we can compose tasks
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

// This is a grouped task, it will allow us to build our own task
// leveraging other tasks
type getHostnameAndIP struct {
	meta *gornir.TaskMetadata
}

// Metadata returns the task metadata
func (t *getHostnameAndIP) Metadata() *gornir.TaskMetadata {
	return t.meta
}

// This is going to be your task result, you can have whatever you want here
type getHostnameAndIPResult struct {
	SubResults []task.RemoteCommandResults // Result of running various commands
}

// If you implement this method on your task result you can control the output when printing it or when using output.RenderResults
func (r getHostnameAndIPResult) String() string {
	return fmt.Sprintf("  - hostname: %s  - ip address: %s", r.SubResults[0].Stdout, r.SubResults[1].Stdout)
}

// Here is where you implement your logic
func (t *getHostnameAndIP) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	// We call the first subtask and store the subresult
	res1, err := (&task.RemoteCommand{Command: "hostname"}).Run(ctx, logger, host)
	if err != nil {
		return getHostnameAndIPResult{}, err
	}

	// We call the second subtask and store the subresult
	res2, err := (&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"}).Run(ctx, logger, host)
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

	// Open an SSH connection towards the devices
	results, err := gr.RunSync(
		&connection.SSHOpen{},
	)
	if err != nil {
		log.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, "Connecting to devices via ssh", true)

	// defer closing the SSH connection we just opened
	defer func() {
		results, err = gr.RunSync(
			&connection.SSHClose{},
		)
		if err != nil {
			log.Fatal(err)
		}
		output.RenderResults(os.Stdout, results, "Close ssh connection", true)
	}()

	// Now we call our "grouped task", which is just a task that uses other tasks
	// In this example we are managing the connection outside the grouped task
	// but we could easily move that inside the grouped task
	results, err = gr.RunSync(
		&getHostnameAndIP{},
	)
	if err != nil {
		log.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, "Let's run a couple of commands", true)
}
