package logger

import (
	"github.com/nornir-automation/gornir/pkg/gornir"
)

// Null is a logger that doesn't do anything. Implements gornir.Logger interface
type Null struct {
}

// NewNull instantiates a new Null logger
func NewNull() *Null {
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
