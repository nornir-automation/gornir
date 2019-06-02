package logger

import (
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	log "github.com/sirupsen/logrus"
)

// Logrus uses github.com/sirupsen/logrus to log messages
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
func (l *Logrus) WithField(field string, value interface{}) gornir.Logger {
	return &Logrus{logger: l.logger.WithFields(log.Fields{field: value})}
}

func (l *Logrus) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logrus) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logrus) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logrus) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logrus) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
