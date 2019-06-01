package task

import (
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/pkg/errors"

	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
)

type RemoteCommand struct {
	Command string
}

type RemoteCommandResults struct {
	Stdout []byte
	Stderr []byte
}

func (r *RemoteCommand) Run(ctx gornir.Context, wg *sync.WaitGroup, jobResult chan *gornir.JobResult) {
	defer wg.Done()
	result := gornir.NewJobResult(ctx)

	sshConfig := &ssh.ClientConfig{
		User: ctx.Host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(ctx.Host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	port := ctx.Host.Port
	if port == 0 {
		port = 22
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ctx.Host.Hostname, port), sshConfig)
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
