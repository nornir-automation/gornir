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

// InventoryPlugin is an Inventory Source
type InventoryPlugin interface {
	Create() (Inventory, error)
}

// New is a Gornir constructor. It is currently no different that new,
// however is a placeholder for any future defaults.
func New() *Gornir {
	return new(Gornir)
}

// SetOption is a funcion that sets one or more options for a given Gornir.
type SetOption func(r *Gornir) error

// Build construct a new Gornir from scratch with options.
func Build(opts ...SetOption) (*Gornir, error) {
	var gornir Gornir
	for _, opt := range opts {
		err := opt(&gornir)
		if err != nil {
			return nil, err
		}
	}
	return &gornir, nil
}

// WithInventory reads the inventory from a file for a Gornir.
func WithInventory(p InventoryPlugin) SetOption {
	return func(g *Gornir) error {
		inv, err := p.Create()
		if err != nil {
			return errors.Wrap(err, "could not read inventory from plugin")
		}
		g.Inventory = &inv
		return nil
	}
}

// WithLogger sets the logging option for a Gornir.
func WithLogger(l Logger) SetOption {
	return func(g *Gornir) error {
		if l == nil {
			errors.New("didn't receive a valid logger")
		}
		g.Logger = l
		return nil
	}
}

// WithFilter provides a FilterFunc to a Gornir to filter the list of hosts.
func WithFilter(f FilterFunc) SetOption {
	return func(g *Gornir) error {
		if f == nil {
			errors.New("didn't receive a valid filter function")
		}
		g.Inventory = g.Inventory.Filter(context.TODO(), f)
		return nil
	}
}

// Build constructs a new Gornir from an existing one.
func (gr Gornir) Build(opts ...SetOption) (*Gornir, error) {
	for _, opt := range opts {
		err := opt(&gr)
		if err != nil {
			return nil, err
		}
	}
	return &gr, nil
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
