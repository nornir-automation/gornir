package task

import (
	"bytes"
	"context"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"

	"github.com/pkg/errors"
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

func (r *RemoteCommand) Run(ctx context.Context, wg *sync.WaitGroup, jp *gornir.JobParameters, jobResult chan *gornir.JobResult) {
	defer wg.Done()
	host := jp.Host()
	result := gornir.NewJobResult(ctx, jp)

	conn, err := host.GetConnection("ssh")
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to retrieve connection"))
		jobResult <- result
		return
	}
	sshConn := conn.(*connection.SSH)

	session, err := sshConn.Client.NewSession()
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to create session"))
		jobResult <- result
		return
	}

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
