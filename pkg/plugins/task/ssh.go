package task

import (
	"bytes"
	"context"
	"fmt"
	"sync"

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

func (r *RemoteCommand) Run(ctx context.Context, wg *sync.WaitGroup, tp *gornir.TaskParameters, jobResult chan *gornir.JobResult) {
	defer wg.Done()
	host := tp.Host()
	result := gornir.NewJobResult(ctx, tp)

	sshConfig := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	port := host.Port
	if port == 0 {
		port = 22
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.Hostname, port), sshConfig)
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to dial"))
		jobResult <- result
		return
	}

	session, err := client.NewSession()
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to create session"))
		jobResult <- result
		return
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(r.Command); err != nil {
		result.SetErr(errors.Wrap(err, "failed to execute command"))
		jobResult <- result
		return
	}
	result.SetData(&RemoteCommandResults{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()})
	jobResult <- result
}
