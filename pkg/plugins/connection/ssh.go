package connection

import (
	"context"
	"fmt"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Client *ssh.Client
}

func (s *SSH) Close() error {
	return s.Client.Close()
}

// String implemente Stringer interface
func (s SSH) String() string {
	if s.Client == nil {
		return "  - connection closed"
	}
	return "  - connection opened"
}

type SSHOpen struct {
}

func (r *SSHOpen) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
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
		return &SSH{}, errors.Wrap(err, "failed to dial")
	}
	host.SetConnection("ssh", &SSH{client})
	return &SSH{client}, nil
}

type SSHClose struct {
}

func (r *SSHClose) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	conn, err := host.GetConnection("ssh")
	if err != nil {
		return &SSH{}, errors.Wrap(err, "failed to retrieve connection")
	}
	sshConn := conn.(*SSH)

	if err := sshConn.Close(); err != nil {
		return &SSH{}, errors.Wrap(err, "failed to close client")
	}
	return &SSH{}, nil
}
