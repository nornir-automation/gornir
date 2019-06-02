package runner

import (
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

func (r ParallelRunner) Run(ctx gornir.Context, task gornir.Task, results chan *gornir.JobResult) error {
	logger := ctx.Logger().WithField("runFunc", getFunctionName(task))
	logger.Debug("starting runner")

	gr := ctx.Gornir()
	if len(gr.Inventory.Hosts) == 0 {
		logger.Warn("no hosts to run against")
		return nil
	}
	r.wg.Add(len(gr.Inventory.Hosts))

	for hostname, host := range gr.Inventory.Hosts {
		logger.WithField("host", hostname).Debug("calling function")
		go task.Run(ctx.ForHost(host), r.wg, results)
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
