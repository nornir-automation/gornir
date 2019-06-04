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

func Parallel() *ParallelRunner {
	return &ParallelRunner{
		wg: &sync.WaitGroup{},
	}
}

func (r ParallelRunner) Run(ctx context.Context, task gornir.Task, hosts map[string]*gornir.Host, tp *gornir.TaskParameters, results chan *gornir.JobResult) error {
	logger := tp.Logger().WithField("runner", "Parallel")
	logger.Debug("starting runner")

	if len(hosts) == 0 {
		logger.Warn("no hosts to run against")
		return nil
	}
	r.wg.Add(len(hosts))

	for hostname, host := range hosts {
		logger.WithField("host", hostname).Debug("calling function")
		go task.Run(ctx, r.wg, tp.ForHost(host), results)
	}
	return nil
}

func (r ParallelRunner) Wait() error {
	r.wg.Wait()
	return nil
}

func (r ParallelRunner) Close() error {
	return nil
}
