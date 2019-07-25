package gornir

import (
	"context"

	"github.com/pkg/errors"
)

type Processor interface {
	TaskStart(context.Context, Logger, Task) error
	TaskCompleted(context.Context, Logger, Task) error
	HostStart(context.Context, Logger, *Host, Task) error
	HostCompleted(context.Context, Logger, *JobResult, *Host, Task) error
}

type Processors []Processor

func (p Processors) TaskStart(ctx context.Context, logger Logger, task Task) error {
	for _, p := range p {
		if err := p.TaskStart(ctx, logger, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'TaskStart'")
		}
	}
	return nil
}

func (p Processors) TaskCompleted(ctx context.Context, logger Logger, task Task) error {
	for _, p := range p {
		if err := p.TaskCompleted(ctx, logger, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'TaskCompleted'")
		}
	}
	return nil
}

func (p Processors) HostStart(ctx context.Context, logger Logger, host *Host, task Task) error {
	for _, p := range p {
		if err := p.HostStart(ctx, logger, host, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'HostStart'")
		}
	}
	return nil
}

func (p Processors) HostCompleted(ctx context.Context, logger Logger, jobResult *JobResult, host *Host, task Task) error {
	for _, p := range p {
		if err := p.HostCompleted(ctx, logger, jobResult, host, task); err != nil {
			return errors.Wrap(err, "problem running processor during 'HostCompleted'")
		}
	}
	return nil
}
