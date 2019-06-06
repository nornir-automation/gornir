// Package gornir implements the core functionality and define the needed interfaces
// to integrate with the framework
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

// Builder defines the steps to construct a Gornir
type Builder interface {
	SetInventory(p InventoryPlugin) Builder
	SetLogger(l Logger) Builder
	SetFilter(f FilterFunc) Builder
	Build() (*Gornir, error)
}

// FromYAMLBuilder is concrete builder
type FromYAMLBuilder struct {
	plugin InventoryPlugin
	logg   Logger
	filter FilterFunc
}

// SetInventory sets the inventory source for the Gornir being built.
func (b *FromYAMLBuilder) SetInventory(p InventoryPlugin) Builder {
	b.plugin = p
	return b
}

// SetLogger sets the logger for the Gornir being built.
func (b *FromYAMLBuilder) SetLogger(l Logger) Builder {
	b.logg = l
	return b
}

// SetFilter applies a host filter to the Gornir being built.
func (b *FromYAMLBuilder) SetFilter(f FilterFunc) Builder {
	b.filter = f
	return b
}

// Build returns a new Gornir with the specified parameters.
func (b *FromYAMLBuilder) Build() (*Gornir, error) {
	gr := new(Gornir)

	if b.plugin != nil {

		inv, err := b.plugin.Create()
		if err != nil {
			return nil, errors.Wrap(err, "could not read inventory from plugin")
		}
		gr.Inventory = &inv
	}

	if b.logg != nil {
		gr.Logger = b.logg
	}

	if b.filter != nil {
		gr.Inventory = gr.Inventory.Filter(context.TODO(), b.filter)
	}
	return gr, nil
}

// NewFromYAML returns a YAML builder
func NewFromYAML() Builder {
	return new(FromYAMLBuilder)
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
