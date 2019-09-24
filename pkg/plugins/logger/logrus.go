package logger

import (
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	log "github.com/sirupsen/logrus"
)

// Logrus uses github.com/sirupsen/logrus to log messages. Implements gornir.Logger Interface
type Logrus struct {
	logger *log.Entry
}

// NewLogrus instantiates a new Logrus logger
func NewLogrus(debug bool) *Logrus {
	logger := &log.Logger{}
	logger.SetFormatter(&log.TextFormatter{})
	if debug {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
	logger.SetOutput(os.Stdout)
	return &Logrus{logger: log.NewEntry(logger)}
}

// NewLogrusFromEntry instantiates a new Logrus logger
func NewLogrusFromEntry(entry *log.Entry) *Logrus {
	return &Logrus{logger: entry}
}

// WithField implements gornir.Logger interface
func (l *Logrus) WithField(field string, value interface{}) gornir.Logger {
	return &Logrus{logger: l.logger.WithFields(log.Fields{field: value})}
}

// Info implements gornir.Logger interface
func (l *Logrus) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Debug implements gornir.Logger interface
func (l *Logrus) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Error implements gornir.Logger interface
func (l *Logrus) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Warn implements gornir.Logger interface
func (l *Logrus) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Fatal implements gornir.Logger interface
func (l *Logrus) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
