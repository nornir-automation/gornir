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
func (r ParallelRunner) Run(ctx context.Context, task gornir.Task, hosts map[string]*gornir.Host, jp *gornir.JobParameters, results chan *gornir.JobResult) error {
	logger := jp.Logger().WithField("runner", "Parallel")
	logger.Debug("starting runner")

	if len(hosts) == 0 {
		logger.Warn("no hosts to run against")
		return nil
	}
	r.wg.Add(len(hosts))

	for hostname, host := range hosts {
		logger.WithField("host", hostname).Debug("calling function")
		go task.Run(ctx, r.wg, jp.ForHost(host), results)
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
