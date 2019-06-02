package runner

import (
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

func (r SortedRunner) Run(ctx gornir.Context, task gornir.Task, results chan *gornir.JobResult) error {
	logger := ctx.Logger().WithField("runFunc", getFunctionName(task))
	logger.Debug("starting runner")

	gr := ctx.Gornir()
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
		host := ctx.Gornir().Inventory.Hosts[hostname]
		logger.WithField("host", hostname).Debug("calling function")
		task.Run(ctx.ForHost(host), wg, results)
	}
	return nil
}

func (r SortedRunner) Wait() error {
	return nil
}

func (r SortedRunner) Close() error {
	return nil
}
