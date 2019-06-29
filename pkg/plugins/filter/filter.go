// Package filter provides a collection of gornir.FilterFunc
package filter

import (
	"github.com/nornir-automation/gornir/pkg/gornir"
)

// WithHostname returns hosts that have a given hostname
func WithHostname(hostname string) gornir.FilterFunc {
	return func(host *gornir.Host) bool {
		return host.Hostname == hostname
	}
}

// Errored returns hosts that returned an error
func Errored(host *gornir.Host) bool {
	return host.Err() != nil
}

// WithError returns hosts that return a given error
func WithError(err error) gornir.FilterFunc {
	return func(host *gornir.Host) bool {
		return host.Err() == err
	}
}

// Not returns the inverse of the given FilterFunc
func Not(filterFunc gornir.FilterFunc) gornir.FilterFunc {
	return func(host *gornir.Host) bool {
		return !filterFunc(host)
	}
}

// And returns if all the filters returned true
func And(filterFuncs ...gornir.FilterFunc) gornir.FilterFunc {
	return func(host *gornir.Host) bool {
		if len(filterFuncs) == 0 {
			return false
		}
		for _, f := range filterFuncs {
			if !f(host) {
				return false
			}
		}
		return true
	}
}

// Or returns if at least a filter returned true
func Or(filterFuncs ...gornir.FilterFunc) gornir.FilterFunc {
	return func(host *gornir.Host) bool {
		for _, f := range filterFuncs {
			if f(host) {
				return true
			}
		}
		return false
	}
}
