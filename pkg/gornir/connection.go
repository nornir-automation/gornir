package gornir

// Connection defines an interface to write connection tasks
type Connection interface {
	Close() error // Close closes the connection
}
