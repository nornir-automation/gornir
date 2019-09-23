// Package gornir implements the core functionality and define the needed
// interfaces to integrate with the framework
package gornir

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Gornir is the main object that glues everything together
type Gornir struct {
	Inventory  *Inventory // Inventory for the object
	Logger     Logger     // Logger for the object
	Runner     Runner     // Runner that will be used to run the task
	Processors Processors // Processors to be used during the execution
	uuid       string     // uuid is a unique identifier used across the logs to match events
}

// New is a Gornir constructor. It is currently no different that new,
// however is a placeholder for any future defaults.
func New() *Gornir {
	return &Gornir{
		Processors: make(Processors, 0),
	}
}

// Clone returns a new instance of Gornir with the same attributes as the receiver
func (gr *Gornir) Clone() *Gornir {
	return &Gornir{
		Inventory:  gr.Inventory,
		Logger:     gr.Logger,
		Runner:     gr.Runner,
		Processors: gr.Processors,
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

// WithProcessors returns a clone of the current Gornir but with the given Processors
func (gr *Gornir) WithProcessors(p Processors) *Gornir {
	c := gr.Clone()
	c.Processors = p
	return c
}

// WithProcessor returns a clone of the current Gornir but with the given Processor appended to the existing ones
func (gr *Gornir) WithProcessor(p Processor) *Gornir {
	c := gr.Clone()
	c.Processors = append(c.Processors, p)
	return c
}

// WithUUID returns a clone of the current Gornir but with the given UUID set. If not
// specifically set gornir will generate one dynamically on each Run
func (gr *Gornir) WithUUID(u string) *Gornir {
	c := gr.Clone()
	c.uuid = u
	return c
}

// UUID returns either the user defined uuid (if set) or a randomized one
func (gr *Gornir) UUID() string {
	if gr.uuid == "" {
		return uuid.New().String()
	}
	return gr.uuid
}

// RunSync will execute the task over the hosts in the inventory using the given runner.
// This function will block until all the tasks are completed.
func (gr *Gornir) RunSync(task Task) (chan *JobResult, error) {
	logger := gr.Logger.WithField("ID", gr.UUID()).WithField("runFunc", getTaskName(task))

	results := make(chan *JobResult, len(gr.Inventory.Hosts))
	defer close(results)

	if err := gr.Processors.TaskStarted(context.Background(), logger, task); err != nil {
		err = errors.Wrap(err, "problem running TaskStart")
		logger.Error(err.Error())
		return results, err
	}

	err := gr.Runner.Run(
		context.Background(),
		logger,
		gr.Processors,
		task,
		gr.Inventory.Hosts,
		results,
	)
	if err != nil {
		err = errors.Wrap(err, "problem calling runner")
		logger.Error(err.Error())
		return results, err
	}
	if err := gr.Runner.Wait(); err != nil {
		err = errors.Wrap(err, "problem waiting for runner")
		logger.Error(err.Error())
		return results, err
	}

	if err := gr.Processors.TaskCompleted(context.Background(), logger, task); err != nil {
		err = errors.Wrap(err, "problem running TaskCompleted")
		logger.Error(err.Error())
		return results, err
	}

	return results, nil
}

// RunAsync will execute the task over the hosts in the inventory using the given runner.
// This function doesn't block, the user can use the method Runnner.Wait instead.
// It's also up to the user to ennsure the channel is closed and that Processors.TaskCompleted is called
func (gr *Gornir) RunAsync(ctx context.Context, task Task, results chan *JobResult) error {
	logger := gr.Logger.WithField("ID", gr.UUID()).WithField("runFunc", getTaskName(task))

	if err := gr.Processors.TaskStarted(ctx, logger, task); err != nil {
		err = errors.Wrap(err, "problem running TaskStart")
		logger.Error(err.Error())
		return err
	}

	err := gr.Runner.Run(
		ctx,
		logger,
		gr.Processors,
		task,
		gr.Inventory.Hosts,
		results,
	)
	if err != nil {
		err = errors.Wrap(err, "problem calling runner")
		logger.Error(err.Error())
		return err
	}

	if err := gr.Processors.TaskCompleted(ctx, logger, task); err != nil {
		err = errors.Wrap(err, "problem running TaskCompleted")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func getTaskName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
