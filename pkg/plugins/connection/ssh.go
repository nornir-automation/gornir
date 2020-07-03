package connection

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/kevinburke/ssh_config"
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

type authMethodResult struct {
	authMethods []ssh.AuthMethod
	err         error
}

func getAuthMethods(host *gornir.Host, logger gornir.Logger) (*[]ssh.AuthMethod, error) {
	var authMethods = []ssh.AuthMethod{ssh.Password(host.Password)}
	// GetStrict should return a default value per `man ssh_config` if this fails, it is because we couldn't parse the config file
	sshPrivKeyFname, err := ssh_config.GetStrict(host.Hostname, "IdentityFile")
	if err != nil {
		return nil, err
	}
	signer, err := getPubKeySigner(host, sshPrivKeyFname, logger)
	// Drop private key auth from the list and fallback to user/pass
	if err != nil {
		return &authMethods, nil
	}
	authMethods = append(authMethods, ssh.PublicKeys(*signer))
	return &authMethods, nil
}

func getPubKeySigner(host *gornir.Host, sshPrivKeyFname string, logger gornir.Logger) (*ssh.Signer, error) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	if sshPrivKeyFname == "~" {
		sshPrivKeyFname = homeDir
	} else if strings.HasPrefix(sshPrivKeyFname, "~/") {
		sshPrivKeyFname = filepath.Join(homeDir, sshPrivKeyFname[2:])
	}
	key, err := ioutil.ReadFile(sshPrivKeyFname)
	if err != nil {
		logger.Debug(fmt.Sprintf("unable to read private key: %v", err))
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logger.Error(fmt.Sprintf("unable to parse private key: %v", err))
		return nil, err
	}
	return &signer, nil
}

// Run implements gornir.Task interface
func (t *SSHOpen) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	authMethods, err := getAuthMethods(host, logger)
	if err != nil {
		return &SSH{}, err
	}
	sshConfig := &ssh.ClientConfig{
		User:            host.Username,
		Auth:            *authMethods,
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
