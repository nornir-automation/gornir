package runner

import (
	"context"
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
func (r ParallelRunner) Run(ctx context.Context, logger gornir.Logger, task gornir.Task, hosts map[string]*gornir.Host, results chan *gornir.JobResult) error {
	logger = logger.WithField("runner", "Parallel")
	logger.Debug("starting runner")

	if len(hosts) == 0 {
		logger.Warn("no hosts to run against")
		return nil
	}
	r.wg.Add(len(hosts))

	for hostname, host := range hosts {
		logger.WithField("host", hostname).Debug("calling function")
		go gornir.TaskWrapper(ctx, logger.WithField("host", hostname), r.wg, task, host, results)
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
