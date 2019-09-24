package runner

import (
	"context"
	"fmt"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

// ParallelRunner will run each task over the hosts in parallel using a goroutines per Host
type ParallelRunner struct {
	wg *sync.WaitGroup
}

// Parallel returns an instantiated ParallelRunner
func Parallel() *ParallelRunner {
	return &ParallelRunner{
		wg: &sync.WaitGroup{},
	}
}

// Run implements the Run method of the gornir.Runner interface
func (r ParallelRunner) Run(ctx context.Context, logger gornir.Logger, processors gornir.Processors, task gornir.Task, hosts map[string]*gornir.Host, results chan *gornir.JobResult) error {
	logger = logger.WithField("runner", "Parallel")
	logger.Debug("starting runner")

	if len(hosts) == 0 {
		logger.Warn("no hosts to run against")
		return nil
	}
	r.wg.Add(len(hosts))

	for hostname, host := range hosts {
		// We need to lock hostname and host or the closure might end up running other hosts
		hostname := hostname
		host := host

		logger.WithField("host", hostname).Debug("calling function")
		go func() {
			if err := gornir.TaskWrapper(ctx, logger.WithField("host", hostname), processors, r.wg, task, host, results); err != nil {
				logger.Error(fmt.Sprintf("problem calling TaskWrapper: %s", err))
			}
		}()
	}
	return nil
}

// Wait implements the Wait method of the gornir.Runner interface
func (r ParallelRunner) Wait() error {
	r.wg.Wait()
	return nil
}

// Close implements the Close method of the gornir.Runner interface
func (r ParallelRunner) Close() error {
	return nil
}
