package connection

import (
	"context"
	"fmt"
	"sync"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Client *ssh.Client
}

type SSHOpen struct {
}

func (s *SSHOpen) Run(ctx context.Context, wg *sync.WaitGroup, jp *gornir.JobParameters, jobResult chan *gornir.JobResult) {
	defer wg.Done()
	host := jp.Host()
	result := gornir.NewJobResult(ctx, jp)

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
		result.SetErr(errors.Wrap(err, "failed to dial"))
		jobResult <- result
		return
	}

	jp.Host().SetConnection("ssh", &SSH{client})
	jobResult <- result
}

type SSHClose struct {
}

func (s *SSHClose) Run(ctx context.Context, wg *sync.WaitGroup, jp *gornir.JobParameters, jobResult chan *gornir.JobResult) {
	defer wg.Done()
	host := jp.Host()
	result := gornir.NewJobResult(ctx, jp)

	conn, err := host.GetConnection("ssh")
	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to retrieve connection"))
		jobResult <- result
		return
	}
	sshConn := conn.(*SSH)

	if err := sshConn.Client.Close(); err != nil {
		result.SetErr(errors.Wrap(err, "failed to close client"))
		jobResult <- result
		return
	}
	jobResult <- result
}
