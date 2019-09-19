// this is the simplest example possible
package main

import (
	"context"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/processor"
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

	rnr := runner.Sorted()

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr).WithProcessor(processor.Render(os.Stdout, true))

	// Open an SSH connection towards the devices
	_, err = gr.RunSync(
		context.Background(),
		&connection.SSHOpen{},
	)
	if err != nil {
		log.Fatal(err)
	}

	// defer closing the SSH connection we just opened
	defer func() {
		_, err = gr.RunSync(
			context.Background(),
			&connection.SSHClose{},
		)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Following call is going to execute the task over all the hosts using the runner.Parallel runner.
	// Said runner is going to handle the parallelization for us. Gornir.RunS is also going to block
	// until the runner has completed executing the task over all the hosts
	_, err = gr.RunSync(
		context.Background(),
		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Now we upload a file. This shows how the ssh connection is shared across tasks of same or different type
	_, err = gr.RunSync(
		context.Background(),
		&task.SFTPUpload{Src: "/etc/hosts", Dst: "/tmp/asd"},
	)
	if err != nil {
		log.Fatal(err)
	}
}
