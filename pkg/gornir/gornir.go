package gornir

import (
	"context"

	"github.com/pkg/errors"
)

type Gornir struct {
	Inventory *Inventory
	Logger    Logger
}

func (g *Gornir) RunS(title string, runner Runner, task Task) (chan *JobResult, error) {
	results := make(chan *JobResult, len(g.Inventory.Hosts))
	err := runner.Run(
		NewContext(context.Background(), title, g, g.Logger),
		task,
		results,
	)
	if err != nil {
		return results, errors.Wrap(err, "problem calling runner")
	}
	if err := runner.Wait(); err != nil {
		return results, errors.Wrap(err, "problem waiting for runner")
	}
	close(results)
	return results, nil
}

func (g *Gornir) RunA(title string, runner Runner, task Task, results chan *JobResult) error {
	err := runner.Run(
		NewContext(context.Background(), title, g, g.Logger),
		task,
		results,
	)
	if err != nil {
		return errors.Wrap(err, "problem calling runner")
	}
	return nil
}
