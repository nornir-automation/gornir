// In this example we can see how we can call the runner asynchronously
// and process the results without having to wait for all the hosts to complete
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
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

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	// The following call will not block
	err = gr.RunAsync(
		context.Background(),
		"What's my hostname?",
		&task.RemoteCommand{Command: "hostname"},
		results,
	)
	if err != nil {
		log.Fatal(err)
	}

	// This goroutine is going to wait for the runner
	// to complete and close the channel when done
	go func() {
		defer close(results)
		rnr.Wait()
	}()

	// While the previous goroutine is waiting for the
	// runner to complete we can start processing
	// the results as they show up
	for {
		select {
		case res, ok := <-results:
			if !ok {
				// channel is closed
				return
			}
			if res.Err() != nil {
				fmt.Printf("ERROR: %s: %s\n", res.JobParameters().Host().Hostname, res.Err().Error())
			} else {
				fmt.Printf("OK: %s: %s\n", res.JobParameters().Host().Hostname, res.Data().(*task.RemoteCommandResults).Stdout)
			}
		case <-time.After(time.Second * 10):
			return
		}
	}

}
