package task

import (
	"bytes"
	"context"
	"fmt"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"

	"github.com/pkg/errors"
)

// RemoteCommand will open a new Session on an already opened ssh connection and execute the given command
type RemoteCommand struct {
	Command string // Command to execute
}

// RemoteCommandResults is the result of calling RemoteCommand
type RemoteCommandResults struct {
	Stdout []byte // Stdout written by the command
	Stderr []byte // Stderr written by the command
}

// String implemente Stringer interface
func (r RemoteCommandResults) String() string {
	return fmt.Sprintf("  - stdout: %s\n  - stderr: %s", r.Stdout, r.Stderr)
}

// Run runs a command on a remote device via ssh
func (r *RemoteCommand) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	conn, err := host.GetConnection("ssh")
	if err != nil {
		return RemoteCommandResults{}, errors.Wrap(err, "failed to retrieve connection")
	}
	sshConn := conn.(*connection.SSH)

	session, err := sshConn.Client.NewSession()
	if err != nil {
		return RemoteCommandResults{}, errors.Wrap(err, "failed to create session")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(r.Command); err != nil {
		return RemoteCommandResults{}, errors.Wrap(err, "failed to execute command")
	}
	return RemoteCommandResults{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, nil
}
