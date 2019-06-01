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
	logger := logger.NewLogrus()

	inventory, err := inventory.FromYAMLFile("hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	results := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))

	rnr := runner.Sequence()
	err = gr.RunA(
		"What's my hostname?",
		rnr,
		&task.RemoteCommand{Command: "hostname"},
		results,
	)
	if err != nil {
		logger.Fatal(err)
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		defer close(results)
		_ = rnr.Wait()
		c <- struct{}{}
	}()

	for {
		select {
		case res := <-results:
			fmt.Println(res)
		case <-c:
			return
		case <-time.After(time.Second * 10):
			return
		}
	}

}
