package runner_test

import (
	"context"
	"sync"
	"time"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

type testTaskSleep struct {
	sleepDuration time.Duration
}

type testTaskSleepResults struct {
	success bool
}

func (t *testTaskSleep) Run(ctx context.Context, wg *sync.WaitGroup, jp *gornir.JobParameters, jr chan *gornir.JobResult) {
	defer wg.Done()
	time.Sleep(t.sleepDuration)
	result := gornir.NewJobResult(ctx, jp)
	result.SetData(&testTaskSleepResults{success: true})
	jr <- result
}

// Null is a logger that doesn't do anything. Implements gornir.Logger interface
type Null struct {
}

// NewNullLogger instantiates a new Null logger
func NewNullLogger() *Null {
	return &Null{}
}

// WithField implements gornir.Logger interface
func (n *Null) WithField(field string, value interface{}) gornir.Logger {
	return n
}

// Info implements gornir.Logger interface
func (n *Null) Info(args ...interface{}) {
}

// Debug implements gornir.Logger interface
func (n *Null) Debug(args ...interface{}) {
}

// Error implements gornir.Logger interface
func (n *Null) Error(args ...interface{}) {
}

// Warn implements gornir.Logger interface
func (n *Null) Warn(args ...interface{}) {
}

// Fatal implements gornir.Logger interface
func (n *Null) Fatal(args ...interface{}) {
}
