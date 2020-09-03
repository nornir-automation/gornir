package runner

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

// SortedRunner will sort the hosts alphabetically and execute the task over
// each host in sequence without any parallelization
type SortedRunner struct {
}

// Sorted returns an instantiated SortedRunner
func Sorted() *SortedRunner {
	return &SortedRunner{}
}

// Run implements the Run method of the gornir.Runner interface
func (r SortedRunner) Run(ctx context.Context, logger gornir.Logger, processors gornir.Processors, task gornir.Task, hosts map[string]*gornir.Host, results chan *gornir.JobResult) error {
	logger = logger.WithField("runner", "Sorted")
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
		if err := gornir.TaskWrapper(ctx, logger.WithField("host", hostname), processors, wg, task, host, results); err != nil {
			logger.Error(fmt.Sprintf("problem calling TaskWrapper: %s", err))
		}
	}
	return nil
}

// Wait implements the Wait method of the gornir.Runner interface
func (r SortedRunner) Wait() error {
	return nil
}

// Close implements the Close method of the gornir.Runner interface
func (r SortedRunner) Close() error {
	return nil
}
