package gornir

import (
	"context"

	"github.com/pkg/errors"
)

// Processor interface implements a set of methods that are called when certain events occur.
// This lets you tap into those events and process them easily as they occur.
//
// To avoid confusions let's clarify the difference between a Task and a TaskInstance:
// - A Task refers to the idea of running a Task over your inventory
// - A TaskInstance refers to a single run of such task over an element of the Inventory.
//
// For instance, if you call `gr.RunSync(MyTask)` and gr has 10 elements you will run
// a single Task and 10 TaskInstances
type Processor interface {
	// TaskStarted is called before a task starts
	TaskStarted(context.Context, Logger, Task) error
	// TaskCompleted is called after all the TaskInstances have completed
	TaskCompleted(context.Context, Logger, Task) error
	// TaskInstanceStarted is called when a TaskInstance starts
	TaskInstanceStarted(context.Context, Logger, *Host, Task) error
	// TaskInstanceStarted is called after a TaskInstance completed
	TaskInstanceCompleted(context.Context, Logger, *JobResult, *Host, Task) error
}

// Processors stores a list of Processor that can be called during gornir's lifetime
// When Procerssors calls the methods of the same name of each Proccesor you have
// to take into account that:
// - It is done sequentially in the order they were defined
// - They are called in a blocking manner
// - If any returns an error the execution is interrupted so the rest won't be called
type Processors []Processor

// TaskStarted calls Processor's methods of the same name
func (p Processors) TaskStarted(ctx context.Context, logger Logger, task Task) error {
	for _, p := range p {
		if err := p.TaskStarted(ctx, logger, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'TaskStart'")
		}
	}
	return nil
}

// TaskCompleted calls Processor's methods of the same name
func (p Processors) TaskCompleted(ctx context.Context, logger Logger, task Task) error {
	for _, p := range p {
		if err := p.TaskCompleted(ctx, logger, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'TaskCompleted'")
		}
	}
	return nil
}

// TaskInstanceStarted calls Processor's methods of the same name
func (p Processors) TaskInstanceStarted(ctx context.Context, logger Logger, host *Host, task Task) error {
	for _, p := range p {
		if err := p.TaskInstanceStarted(ctx, logger, host, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'HostStart'")
		}
	}
	return nil
}

// TaskInstanceCompleted calls Processor's methods of the same name
func (p Processors) TaskInstanceCompleted(ctx context.Context, logger Logger, jobResult *JobResult, host *Host, task Task) error {
	for _, p := range p {
		if err := p.TaskInstanceCompleted(ctx, logger, jobResult, host, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'HostCompleted'")
		}
	}
	return nil
}
