package gornir

import (
	"github.com/pkg/errors"
)

// Host represent a host
type Host struct {
	err         error
	Port        uint16                 `yaml:"port"`     // Port to connect to
	Hostname    string                 `yaml:"hostname"` // Hostname/FQDN/IP to connect to
	Username    string                 `yaml:"username"` // Username to use for authentication purposes
	Password    string                 `yaml:"password"` // Password to use for authentication purposes
	Platform    string                 `yaml:"platform"` // Platform of the device
	Data        map[string]interface{} `yaml:"data"`     // Data belonging to the host
	connections map[string]Connection
}

// Inventory represents a collection of Hosts
type Inventory struct {
	Hosts map[string]*Host // Hosts represents a collection of Hosts
}

// FilterFunc is a function that can be used to filter the inventory
type FilterFunc func(*Host) bool

// Filter filters the hosts in the inventory returning a copy of the current
// Inventory instance but with only the hosts that passed the filter
func (i *Inventory) Filter(f FilterFunc) *Inventory {
	filtered := &Inventory{
		Hosts: make(map[string]*Host),
	}
	for hostname, host := range i.Hosts {
		if f(host) {
			filtered.Hosts[hostname] = host
		}
	}
	return filtered
}

// SetErr stores the error in the host
func (h *Host) SetErr(err error) {
	h.err = err
}

// Err returns the stored error
func (h *Host) Err() error {
	return h.err
}

// SetConnections stores a connection
func (h *Host) SetConnection(name string, conn Connection) {
	if h.connections == nil {
		h.connections = make(map[string]Connection)
	}
	h.connections[name] = conn
}

// GetConnection retrieves a connection that was previously set
func (h *Host) GetConnection(name string) (Connection, error) {
	if h.connections == nil {
		h.connections = make(map[string]Connection)
	}
	if c, ok := h.connections[name]; ok {
		return c, nil
	}
	return nil, errors.New("couldn't find connection")
}
