package connection

import (
	"context"
	"testing"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"golang.org/x/crypto/ssh"
)

func TestSSHOpen(t *testing.T) {
	var haveIBeenCalled bool
	cfgfn := func(host *gornir.Host, logger gornir.Logger) (*ssh.ClientConfig, error) {
		haveIBeenCalled = true
		return &ssh.ClientConfig{}, nil
	}

	testCases := []struct {
		name         string
		fn           ClientConfigFn
		expectCalled bool
	}{
		{"Running with defaults", nil, false},
		{"Running with a custom client config", cfgfn, true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			haveIBeenCalled = false
			runner := &SSHOpen{ClientConfigFn: tc.fn}
			// ignore the error returned from `Run` since we expect the SSH dial to fail
			runner.Run(context.Background(), logger.NewLogrus(false), &gornir.Host{Hostname: "dev1"}) // nolint
			if haveIBeenCalled != tc.expectCalled {
				t.Errorf("got %t; want %t", haveIBeenCalled, tc.expectCalled)
			}
		})
	}
}
