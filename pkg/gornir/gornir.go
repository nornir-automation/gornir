// Package Gornir implements the core functionality and define the needed interfaces to integrate with the framework
package gornir

import (
	"context"

	"github.com/pkg/errors"
)

// Gornir is the main object that glues everything together
type Gornir struct {
	Inventory *Inventory // Inventory for the object
	Logger    Logger     // Logger for the object
}

// Filter filters the hosts in the inventory returning a copy of the current
// Gornir instance but with only the hosts that passed the filter
func (g *Gornir) Filter(f FilterFunc) *Gornir {
	return &Gornir{
		Inventory: g.Inventory.Filter(g, f),
		Logger:    g.Logger,
	}
}

// RunS will execute the task over the hosts in the inventory using the given runner.
// This function will block until all the tasks are completed.
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

// RunA will execute the task over the hosts in the inventory using the given runner.
// This function doesn't block, the user can use the method Runnner.Wait instead.
// It's also up to the user to ennsure the channel is closed
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
