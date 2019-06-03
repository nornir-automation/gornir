package gornir

import (
	"runtime"

	"github.com/google/uuid"
)

// Params implements the context.Params interface and enriches it
// with extra useful information. You will find this object mosty
// in two places; in the JobResult and in objects implementing
// the Task interface
type Params struct {
	title  string
	id     string
	gr     *Gornir
	logger Logger
	host   *Host
}

// NewParams returns a new Params
func NewParams(title string, gr *Gornir, logger Logger) Params {
	id := uuid.New().String()
	return Params{
		id:     id,
		gr:     gr,
		logger: logger.WithField("ID", id),
		title:  title,
	}
}

// ForHost returns a copy of the Params adding the Host to it
func (c Params) ForHost(host *Host) Params {
	return Params{
		id:     c.id,
		gr:     c.gr,
		logger: c.logger,
		title:  c.title,
		host:   host,
	}
}

// Title returns the title of the task
func (c Params) Title() string {
	return c.title
}

// Gornir returns the Gornir object that triggered the execution of the task
func (c Params) Gornir() *Gornir {
	return c.gr
}

// Host returns the Host associated with the context
func (c Params) Host() *Host {
	return c.host
}

// ID returns the unique ID associated with the execution. All hosts
// will share the same ID for a given Run
func (c Params) ID() string {
	return c.id
}

// Logger returns a ready to use Logger
func (c Params) Logger() Logger {
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
