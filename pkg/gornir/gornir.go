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
func (gr *Gornir) Filter(ctx context.Context, f FilterFunc) *Gornir {
	return &Gornir{
		Inventory: gr.Inventory.Filter(ctx, f),
		Logger:    gr.Logger,
	}
}

// RunSync will execute the task over the hosts in the inventory using the given runner.
// This function will block until all the tasks are completed.
func (gr *Gornir) RunSync(title string, runner Runner, task Task) (chan *JobResult, error) {
	results := make(chan *JobResult, len(gr.Inventory.Hosts))
	err := runner.Run(
		context.Background(),
		task,
		&TaskParameters{
			Title:  title,
			Gornir: gr,
			Logger: gr.Logger,
		},
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

// RunAsync will execute the task over the hosts in the inventory using the given runner.
// This function doesn't block, the user can use the method Runnner.Wait instead.
// It's also up to the user to ennsure the channel is closed
func (gr *Gornir) RunAsync(ctx context.Context, title string, runner Runner, task Task, results chan *JobResult) error {
	err := runner.Run(
		context.Background(), // TODO pass this?
		task,
		&TaskParameters{
			Title:  title,
			Gornir: gr,
			Logger: gr.Logger,
		},
		results,
	)
	if err != nil {
		return errors.Wrap(err, "problem calling runner")
	}
	return nil
}
