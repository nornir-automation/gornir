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
}

// New is a Gornir constructor. It is currently no different that new,
// however is a placeholder for any future defaults.
func New() *Gornir {
	return new(Gornir)
}

// InventoryPlugin is an Inventory Source
type InventoryPlugin interface {
	Create() (Inventory, error)
}

// WithInventory creates a new Gornir with an Inventory.
func (gr *Gornir) WithInventory(inv Inventory) *Gornir {
	return &Gornir{
		Logger:    gr.Logger,
		Inventory: &inv,
	}
}

// WithFilter creates a new Gornir with a filtered Inventory.
func (gr *Gornir) WithFilter(f FilterFunc) *Gornir {
	return &Gornir{
		Logger:    gr.Logger,
		Inventory: gr.Inventory.Filter(context.TODO(), f),
	}
}

// WithLogger creates a new Gornir with a Logger.
func (gr *Gornir) WithLogger(l Logger) *Gornir {
	return &Gornir{
		Logger:    l,
		Inventory: gr.Inventory,
	}
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
	logger := gr.Logger.WithField("ID", uuid.New().String()).WithField("runFunc", getFunctionName(task))
	results := make(chan *JobResult, len(gr.Inventory.Hosts))
	defer close(results)
	err := runner.Run(
		context.Background(),
		task,
		gr.Inventory.Hosts,
		NewJobParameters(title, logger),
		results,
	)
	if err != nil {
		return results, errors.Wrap(err, "problem calling runner")
	}
	if err := runner.Wait(); err != nil {
		return results, errors.Wrap(err, "problem waiting for runner")
	}
	return results, nil
}

// RunAsync will execute the task over the hosts in the inventory using the given runner.
// This function doesn't block, the user can use the method Runnner.Wait instead.
// It's also up to the user to ennsure the channel is closed
func (gr *Gornir) RunAsync(ctx context.Context, title string, runner Runner, task Task, results chan *JobResult) error {
	logger := gr.Logger.WithField("ID", uuid.New().String()).WithField("runFunc", getFunctionName(task))
	err := runner.Run(
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
