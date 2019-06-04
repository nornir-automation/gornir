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

func (r SortedRunner) Run(ctx context.Context, task gornir.Task, hosts map[string]*gornir.Host, jp *gornir.JobParameters, results chan *gornir.JobResult) error {
	logger := jp.Logger().WithField("runner", "Sorted")
	logger.Debug("starting runner")

	if len(hosts) == 0 {
		logger.Warn("no hosts to run against")
		return nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(hosts))

	sortedHostnames := make([]string, len(hosts))
	i := 0
	for hostname := range hosts {
		sortedHostnames[i] = hostname
		i++
	}
	sort.Slice(sortedHostnames, func(i, j int) bool { return sortedHostnames[i] < sortedHostnames[j] })

	for _, hostname := range sortedHostnames {
		host := hosts[hostname]
		logger.WithField("host", hostname).Debug("calling function")
		task.Run(ctx, wg, jp.ForHost(host), results)
	}
	return nil
}

func (r SortedRunner) Wait() error {
	return nil
}

func (r SortedRunner) Close() error {
	return nil
}
