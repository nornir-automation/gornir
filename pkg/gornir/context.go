package gornir

import (
	"context"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// Context implements the context.Context interface and enriches it
// with extra useful information. You will find this object mosty
// in two places; in the JobResult and in objects implementing
// the Task interface
type Context struct {
	ctx    context.Context
	title  string
	id     string
	gr     *Gornir
	logger Logger
	host   *Host
}

// NewContext returns a new Context
func NewContext(ctx context.Context, title string, gr *Gornir, logger Logger) Context {
	id := uuid.New().String()
	return Context{
		id:     id,
		gr:     gr,
		ctx:    ctx,
		logger: logger.WithField("ID", id),
		title:  title,
	}
}

// ForHost returns a copy of the Context adding the Host to it
func (c Context) ForHost(host *Host) Context {
	return Context{
		id:     c.id,
		gr:     c.gr,
		ctx:    c.ctx,
		logger: c.logger,
		title:  c.title,
		host:   host,
	}
}

// Title returns the title of the task
func (c Context) Title() string {
	return c.title
}

// Gornir returns the Gornir object that triggered the execution of the task
func (c Context) Gornir() *Gornir {
	return c.gr
}

// Host returns the Host associated with the context
func (c Context) Host() *Host {
	return c.host
}

// Deadline delegates the method to the underlying context pass upon creation
func (c Context) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

// Done delegates the method to the underlying context pass upon creation
func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Value delegates the method to the underlying context pass upon creation
func (c Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// Err will return the error returned by a task. Otherwise it will be nil
func (c Context) Err() error {
	return c.ctx.Err()
}

// ID returns the unique ID associated with the execution. All hosts
// will share the same ID for a given Run
func (c Context) ID() string {
	return c.id
}

// Logger returns a ready to use Logger
func (c Context) Logger() Logger {
	return c.logger.WithField("ID", c.id).WithField("funcName", getFrame(1).Function)
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}
	return frame
}
