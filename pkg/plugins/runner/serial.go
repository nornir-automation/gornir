package runner

import (
	"context"
	"sort"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

// SortedRunner will sort the hosts alphabetically and execute the task over
// each host in sequence without any parallelization
type SortedRunner struct {
}

func Sorted() *SortedRunner {
	return &SortedRunner{}
}

func (r SortedRunner) Run(ctx context.Context, task gornir.Task, taskParameters *gornir.TaskParameters, results chan *gornir.JobResult) error {
	logger := taskParameters.Logger.WithField("runFunc", getFunctionName(task))
	logger.Debug("starting runner")

	gr := taskParameters.Gornir
	if len(gr.Inventory.Hosts) == 0 {
		logger.Warn("no hosts to run against")
		return nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(gr.Inventory.Hosts))

	sortedHostnames := make([]string, len(gr.Inventory.Hosts))
	i := 0
	for hostname := range gr.Inventory.Hosts {
		sortedHostnames[i] = hostname
		i++
	}
	sort.Slice(sortedHostnames, func(i, j int) bool { return sortedHostnames[i] < sortedHostnames[j] })

	for _, hostname := range sortedHostnames {
		host := gr.Inventory.Hosts[hostname]
		logger.WithField("host", hostname).Debug("calling function")
		task.Run(ctx, wg, taskParameters.ForHost(host), results)
	}
	return nil
}

func (r SortedRunner) Wait() error {
	return nil
}

func (r SortedRunner) Close() error {
	return nil
}
