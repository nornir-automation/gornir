package runner

import (
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

type SequenceParams struct {
	wg *sync.WaitGroup
}

func Sequence() *SequenceParams {
	return &SequenceParams{
		wg: &sync.WaitGroup{},
	}
}

func (r SequenceParams) Run(ctx gornir.Context, task gornir.Task, results chan *gornir.JobResult) error {
	logger := ctx.Logger().WithField("runFunc", gornir.GetFunctionName(task))
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

func (r SequenceParams) Wait() error {
	r.wg.Wait()
	return nil
}

func (r SequenceParams) Close() error {
	return nil
}
