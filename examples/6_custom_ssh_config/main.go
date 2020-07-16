// this is the simplest example possible
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
	"golang.org/x/crypto/ssh"
)

func getPubKeySigner(host *gornir.Host, sshPrivKeyFname string, logger gornir.Logger) (*ssh.Signer, error) {
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

func GetSSHConfig(host *gornir.Host, logger gornir.Logger) (*ssh.ClientConfig, error) {
	var authMethods = []ssh.AuthMethod{ssh.Password(host.Password)}
	// Under normal circumstances, you probably want to use something like the github.com/kevinburke/ssh_config package
	// usr, _ := user.Current()
	// homeDir := usr.HomeDir
	// if sshPrivKeyFname == "~" {
	//     sshPrivKeyFname = homeDir
	// } else if strings.HasPrefix(sshPrivKeyFname, "~/") {
	//     sshPrivKeyFname = filepath.Join(homeDir, sshPrivKeyFname[2:])
	// }
	// sshPrivKeyFname, err := ssh_config.GetStrict(host.Hostname, "IdentityFile")
	// GetStrict should return a default value per `man ssh_config` if this fails, it is because we couldn't parse the config file
	// if err != nil {
	// 	return nil, err
	// }
	sshPrivKeyFname := "/go/src/github.com/nornir-automation/gornir/examples/6_custom_ssh_config/id_rsa"
	signer, err := getPubKeySigner(host, sshPrivKeyFname, logger)
	if err != nil {
		return nil, err
	}
	authMethods = append(authMethods, ssh.PublicKeys(*signer))
	sshConfig := &ssh.ClientConfig{
		User:            host.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	} // #nosec
	return sshConfig, nil
}

func main() {
	// Instantiate a logger plugin
	log := logger.NewLogrus(false)

	// Load the inventory using the FromYAMLFile plugin
	file := "/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml"
	plugin := inventory.FromYAML{HostsFile: file}
	inv, err := plugin.Create()
	if err != nil {
		log.Fatal(err)
	}

	rnr := runner.Sorted()

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	// Open an SSH connection towards the devices
	results, err := gr.RunSync(
		context.Background(),
		&connection.SSHOpen{ClientConfigFn: GetSSHConfig},
	)
	if err != nil {
		log.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, "Connecting to devices via ssh", true)

	// defer closing the SSH connection we just opened
	defer func() {
		results, err = gr.RunSync(
			context.Background(),
			&connection.SSHClose{},
		)
		if err != nil {
			log.Fatal(err)
		}
		output.RenderResults(os.Stdout, results, "Close ssh connection", true)
	}()

	// Following call is going to execute the task over all the hosts using the runner.Parallel runner.
	// Said runner is going to handle the parallelization for us. Gornir.RunS is also going to block
	// until the runner has completed executing the task over all the hosts
	results, err = gr.RunSync(
		context.Background(),
		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
	)
	if err != nil {
		log.Fatal(err)
	}
	// next call is going to print the result on screen
	output.RenderResults(os.Stdout, results, "What is my ip?", true)

	// Now we upload a file. This shows how the ssh connection is shared across tasks of same or different type
	results, err = gr.RunSync(
		context.Background(),
		&task.SFTPUpload{Src: "/etc/hosts", Dst: "/tmp/asd"},
	)
	if err != nil {
		log.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, "Upload File", true)
}
