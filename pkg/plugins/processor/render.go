package processor

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

const (
	redColor   = "\u001b[31m"
	greenColor = "\u001b[32m"
	blueColor  = "\u001b[34m"
	resetColor = "\u001b[0m"
)

func red(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", redColor, m, resetColor)
	}
	return m
}

func green(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", greenColor, m, resetColor)
	}
	return m
}

func blue(m string, color bool) string {
	if color {
		return fmt.Sprintf("%v%v%v", blueColor, m, resetColor)
	}
	return m
}

// RenderProcessor is a processor that writes the result to an io.Writer
type RenderProcessor struct {
	mux   *sync.Mutex
	wr    io.Writer
	color bool
}

// Render returns a configured RenderProcessor
func Render(wr io.Writer, color bool) *RenderProcessor {
	return &RenderProcessor{
		mux:   &sync.Mutex{},
		wr:    wr,
		color: color,
	}
}

// TaskStarted renders task.Metdata().Identifier or the task's struct's name
func (r *RenderProcessor) TaskStarted(ctx context.Context, logger gornir.Logger, task gornir.Task) error {
	var taskName string

	meta := task.Metadata()
	if meta != nil {
		taskName = meta.Identifier
	}

	if taskName == "" {
		if t := reflect.TypeOf(task); t.Kind() == reflect.Ptr {
			taskName = t.Elem().Name()
		} else {
			taskName = t.Name()
		}
	}

	_, err := r.wr.Write([]byte(blue(fmt.Sprintf("# %s\n", taskName), r.color)))
	return err
}

// TaskCompleted doesn't do anything
func (r *RenderProcessor) TaskCompleted(ctx context.Context, logger gornir.Logger, task gornir.Task) error {
	return nil
}

// TaskInstanceStarted doesn't do anything
func (r *RenderProcessor) TaskInstanceStarted(ctx context.Context, logger gornir.Logger, host *gornir.Host, task gornir.Task) error {
	return nil
}

// TaskInstanceCompleted renders either the result or the error resulted in the execution of the TaskInstance
func (r *RenderProcessor) TaskInstanceCompleted(ctx context.Context, logger gornir.Logger, jobResult *gornir.JobResult, host *gornir.Host, task gornir.Task) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	switch jobResult.Err() {
	case nil:
		if _, err := r.wr.Write([]byte(green(fmt.Sprintf("@ %s\n", host.Hostname), r.color))); err != nil {
			return err
		}
		if _, err := r.wr.Write([]byte(fmt.Sprintf("%v\n\n", jobResult.Data()))); err != nil {
			return err
		}
	default:
		if _, err := r.wr.Write([]byte(red(fmt.Sprintf("@ %s\n", host.Hostname), r.color))); err != nil {
			return err
		}
		if _, err := r.wr.Write([]byte(fmt.Sprintf("  - err: %v\n\n", jobResult.Err()))); err != nil {
			return err
		}
	}
	return nil
}
