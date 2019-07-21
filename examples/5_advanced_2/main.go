// In this example we can see how we can call the runner asynchronously
// and process the results without having to wait for all the hosts to complete
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

// This is a grouped task, it will allow us to build our own task
// leveraging other tasks
type getHostnameAndIP struct {
}

// This is going to be your task result, you can have whatever you want here
type getHostnameAndIPResult struct {
	SubResults []gornir.TaskInstanceResult // Result of running various commands
}

func (r getHostnameAndIPResult) String() string {
	res := ""
	for _, r := range r.SubResults {
		res += fmt.Sprintf("%s\n", r)
	}
	return res
}

// Here is where you implement your logic
func (r *getHostnameAndIP) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	resOpen, err := (&connection.SSHOpen{}).Run(ctx, logger, host)
	if err != nil {
		return getHostnameAndIPResult{}, err
	}

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

	resClose, err := (&connection.SSHClose{}).Run(ctx, logger, host)
	if err != nil {
		return getHostnameAndIPResult{}, err
	}

	return getHostnameAndIPResult{
		SubResults: []gornir.TaskInstanceResult{
			resOpen,
			res1,
			res2,
			resClose,
		},
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

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	// The following call will not block
	err = gr.RunAsync(
		context.Background(),
		&getHostnameAndIP{},
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
				fmt.Printf("ERROR: %s: %s\n", res.Host().Hostname, res.Err().Error())
			} else {
				fmt.Printf("OK: %s:\n%s\n", res.Host().Hostname, res.Data().(gornir.TaskInstanceResult))
			}
		case <-time.After(time.Second * 10):
			return
		}
	}

}
