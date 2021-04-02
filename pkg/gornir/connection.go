package gornir

import (
	"context"
)

// Connection defines an interface to write connection tasks
type Connection interface {
	Close(context.Context) error // Close closes the connection
}
