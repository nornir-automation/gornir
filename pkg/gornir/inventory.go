package gornir

import (
	"context"
)

// Host represent a host
type Host struct {
	err      error
	Port     uint8  `yaml:"port"`     // Port to connect to
	Hostname string `yaml:"hostname"` // Hostname/FQDN/IP to connect to
	Username string `yaml:"username"` // Username to use for authentication purposes
	Password string `yaml:"password"` // Password to use for authentication purposes
	Platform string `yaml:"platform"` // Platform of the device
}

// Inventory represents a collection of Hosts
type Inventory struct {
	Hosts map[string]*Host // Hosts represents a collection of Hosts
}

// FilterFunc is a function that can be used to filter the inventory
type FilterFunc func(context.Context, *Host) bool

// Filter filters the hosts in the inventory returning a copy of the current
// Inventory instance but with only the hosts that passed the filter
func (i *Inventory) Filter(ctx context.Context, f FilterFunc) *Inventory {
	filtered := &Inventory{
		Hosts: make(map[string]*Host),
	}
	for hostname, host := range i.Hosts {
		if f(ctx, host) {
			filtered.Hosts[hostname] = host
		}
	}
	return filtered
}

func (h *Host) setErr(err error) {
	h.err = err
}

// Err returns the stored error
func (h *Host) Err() error {
	return h.err
}
