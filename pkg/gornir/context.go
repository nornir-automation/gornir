package gornir

import (
	"context"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type Context struct {
	ctx    context.Context
	title  string
	id     string
	gr     *Gornir
	logger Logger
	Host   *Host
}

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

func (c Context) ForHost(host *Host) Context {
	return Context{
		id:     c.id,
		gr:     c.gr,
		ctx:    c.ctx,
		logger: c.logger,
		title:  c.title,
		Host:   host,
	}
}

func (c Context) Title() string {
	return c.title
}

func (c Context) Gornir() *Gornir {
	return c.gr
}

func (c Context) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c Context) Err() error {
	return c.ctx.Err()
}

func (c Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c Context) ID() string {
	return c.id
}

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
