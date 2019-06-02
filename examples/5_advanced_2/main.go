// In this example we can see how we can call the runner asynchronously
// and process the results without having to wait for all the hosts to complete
package main

import (
	"fmt"
	"time"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	logger := logger.NewLogrus(false)

	inventory, err := inventory.FromYAMLFile("/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	rnr := runner.Parallel()

	// The following call will not block
	err = gr.RunA(
		"What's my hostname?",
		rnr,
		&task.RemoteCommand{Command: "hostname"},
		results,
	)
	if err != nil {
		logger.Fatal(err)
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
			fmt.Println(res)
		case <-time.After(time.Second * 10):
			return
		}
	}

}
