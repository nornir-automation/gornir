package task

import (
	"bytes"
	"context"
	"fmt"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// RemoteCommand will connect to the Host via ssh and execute the given command
type RemoteCommand struct {
	Command string // Command to execute
}

// RemoteCommandResults will be accessible via JobResult.Data()
type RemoteCommandResults struct {
	Stdout []byte // Stdout written by the command
	Stderr []byte // Stderr written by the command
}

func (r RemoteCommandResults) String() string {
	return fmt.Sprintf("    stdout: %s\n    stderr: %s", r.Stdout, r.Stderr)
}

func (r *RemoteCommand) Run(ctx context.Context, host *gornir.Host) (interface{}, error) {
	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	} // #nosec
	port := host.Port
	if port == 0 {
		port = 22
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.Hostname, port), sshConfig)
	if err != nil {
		return RemoteCommandResults{}, errors.Wrap(err, "failed to dial")
	}

	session, err := client.NewSession()
	if err != nil {
		return RemoteCommandResults{}, errors.Wrap(err, "failed to create session")
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(r.Command); err != nil {
		return RemoteCommandResults{}, errors.Wrap(err, "failed to execute command")
	}
	return RemoteCommandResults{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, nil
}
