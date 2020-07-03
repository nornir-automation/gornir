package connection

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// SSH is a Connection plugins that connects to device via ss using the golang.org/x/crypto/ssh
// package. Current implementation only supports authentication with a password and has
// ssh.InsecureIgnoreHostKey set
type SSH struct {
	Client *ssh.Client
}

// Close closes the connection
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

// SSHOpen is a Connection plugin that opens a connection with a device
type SSHOpen struct {
	Meta *gornir.TaskMetadata // Task metadata
}

// Metadata returns the task metadata
func (t *SSHOpen) Metadata() *gornir.TaskMetadata {
	return t.Meta
}

func configPubKey(host *gornir.Host) (*ssh.Signer, error) {
	// var hostKey ssh.PublicKey
	sshPrivKeyFname, ok := host.Data["ssh_private_key_file"].(string)
	if !ok {
		err := errors.New("'ssh_private_key_file' is null")
		log.Printf("%v", err)
		return nil, err
	}
	key, err := ioutil.ReadFile(sshPrivKeyFname)
	if err != nil {
		log.Printf("unable to read private key: %v", err)
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Printf("unable to parse private key: %v", err)
		return nil, err
	}
	return &signer, nil
}

// Run implements gornir.Task interface
func (t *SSHOpen) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	var authMethods = []ssh.AuthMethod{ssh.Password(host.Password)}
	signer, err := configPubKey(host)
	if err == nil {
		authMethods = append(authMethods, ssh.PublicKeys(*signer))
	}
	sshConfig := &ssh.ClientConfig{
		User:            host.Username,
		Auth:            authMethods,
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

// SSHClose is a Connection plugin that closes an already opened ssh connection
type SSHClose struct {
	Meta *gornir.TaskMetadata // Task metadata
}

// Metadata returns the task metadata
func (t *SSHClose) Metadata() *gornir.TaskMetadata {
	return t.Meta
}

// Run implements gornir.Task interface
func (t *SSHClose) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
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
