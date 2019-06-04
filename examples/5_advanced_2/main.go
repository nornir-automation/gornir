// In this example we can see how we can call the runner asynchronously
// and process the results without having to wait for all the hosts to complete
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	// Instantiate a logger plugin.
	logger := logger.NewLogrus(false) 
	// File where the inventory will be loaded from.
	file := "/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml"
	
	// Instantiate Gornir
	gr, err := gornir.Build(
		gornir.WithInventory(file),
		gornir.WithLogger(logger),
	)	
	if err != nil {
		logger.Fatal(err)
	}

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	rnr := runner.Parallel()

	// The following call will not block
	err = gr.RunAsync(
		context.Background(),
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
