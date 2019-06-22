package connections

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

var ciphers = []string{"3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc", "aes128-ctr"}

type SSHConn struct {
	Host     string
	Port     uint8
	Username string
	Password string
	session  *ssh.Session
}

func NewSSHConn(host string, port uint8, username string, password string) *SSHConn {
	return &SSHConn{
		Host:     host,
		Username: username,
		Password: password,
		Port:     port,
		session:  nil,
	}
}

func (c *SSHConn) Open() error {

	sshConfig := &ssh.ClientConfig{
		User:            c.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         6 * time.Second}

	sshConfig.Ciphers = append(sshConfig.Ciphers, ciphers...)
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return err
	}

	session, err := conn.NewSession()

	if err != nil {
		return err
	}
	c.session = session
	return nil
}

func (c *SSHConn) Close() {
	c.session.Close()
}

func (c *SSHConn) Send(cmd string) (RemoteCommandResults, error) {
	var results RemoteCommandResults
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c.session.Stdout = &stdout
	c.session.Stderr = &stderr

	err := c.session.Run(cmd)
	if err != nil {
		return results, err
	}

	results.Stdout = stdout.String()
	results.Stderr = stderr.String()

	return results, nil

}
