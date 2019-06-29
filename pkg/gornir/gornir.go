// Package gornir implements the core functionality and define the needed
// interfaces to integrate with the framework
package gornir

import (
	"context"
	"reflect"
	"runtime"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Gornir is the main object that glues everything together
type Gornir struct {
	Inventory *Inventory // Inventory for the object
	Logger    Logger     // Logger for the object
	Runner    Runner     // Runner that will be used to run the task
}

// New is a Gornir constructor. It is currently no different that new,
// however is a placeholder for any future defaults.
func New() *Gornir {
	return new(Gornir)
}

// Clone returns a new instance of Gornir with the same attributes as the receiver
func (gr *Gornir) Clone() *Gornir {
	return &Gornir{
		Inventory: gr.Inventory,
		Logger:    gr.Logger,
		Runner:    gr.Runner,
	}
}

// WithRunner returns a clone of the current Gornir but with the given runner
func (gr *Gornir) WithRunner(rnr Runner) *Gornir {
	c := gr.Clone()
	c.Runner = rnr
	return c
}

// WithInventory returns a clone of the current Gornir but with the given inventory
func (gr *Gornir) WithInventory(inv Inventory) *Gornir {
	c := gr.Clone()
	c.Inventory = &inv
	return c
}

// Filter creates a new Gornir with a filtered Inventory.
// It filters the hosts in the inventory returning a copy of the current
// Gornir instance but with only the hosts that passed the filter.
func (gr *Gornir) Filter(f FilterFunc) *Gornir {
	c := gr.Clone()
	c.Inventory = c.Inventory.Filter(f)
	return c
}

// WithLogger returns a clone of the current Gornir but with the given logger
func (gr *Gornir) WithLogger(l Logger) *Gornir {
	c := gr.Clone()
	c.Logger = l
	return c
}

// RunSync will execute the task over the hosts in the inventory using the given runner.
// This function will block until all the tasks are completed.
func (gr *Gornir) RunSync(title string, task Task) (chan *JobResult, error) {
	logger := gr.Logger.WithField("ID", uuid.New().String()).WithField("runFunc", getFunctionName(task))
	results := make(chan *JobResult, len(gr.Inventory.Hosts))
	defer close(results)
	err := gr.Runner.Run(
		context.Background(),
		task,
		gr.Inventory.Hosts,
		NewJobParameters(title, logger),
		results,
	)
	if err != nil {
		return results, errors.Wrap(err, "problem calling runner")
	}
	if err := gr.Runner.Wait(); err != nil {
		return results, errors.Wrap(err, "problem waiting for runner")
	}
	return results, nil
}

// RunAsync will execute the task over the hosts in the inventory using the given runner.
// This function doesn't block, the user can use the method Runnner.Wait instead.
// It's also up to the user to ennsure the channel is closed
func (gr *Gornir) RunAsync(ctx context.Context, title string, task Task, results chan *JobResult) error {
	logger := gr.Logger.WithField("ID", uuid.New().String()).WithField("runFunc", getFunctionName(task))
	err := gr.Runner.Run(
		ctx,
		task,
		gr.Inventory.Hosts,
		NewJobParameters(title, logger),
		results,
	)
	if err != nil {
		return errors.Wrap(err, "problem calling runner")
	}
	return nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
