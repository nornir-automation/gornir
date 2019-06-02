package gornir

// Logger defines the interface that a logger plugin can implement to
// provide logging capabilities
type Logger interface {
	Info(...interface{})                  // Info logs an informational message
	Debug(...interface{})                 // Debug logs a debug message
	Error(...interface{})                 // Error logs an error
	Warn(...interface{})                  // Warn logs a warning
	Fatal(...interface{})                 // Fatal logs a fatal event
	WithField(string, interface{}) Logger // WithField adds data to the logger
}
