// Package task implements various Task plugins that can be run over Hosts
package task

import (
	"context"
	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connections"
	"github.com/pkg/errors"
	"sync"
)

// RemoteCommand will connect to the Host via ssh and execute the given command
type RemoteCommand struct {
	Command string // Command to execute
}

func (r *RemoteCommand) Run(ctx context.Context, wg *sync.WaitGroup, jp *gornir.JobParameters, jobResult chan *gornir.JobResult) {
	defer wg.Done()
	host := jp.Host()
	result := gornir.NewJobResult(ctx, jp)

	port := host.Port
	if port == 0 {
		port = 22
	}

	conn := connections.NewSSHConn(host.Hostname, port, host.Username, host.Password)
	if err := conn.Open(); err != nil {
		result.SetErr(errors.Wrap(err, "failed to Open Device"))
		jobResult <- result
		return
	}
	defer conn.Close()

	remoteResult, err := conn.Send(r.Command)

	if err != nil {
		result.SetErr(errors.Wrap(err, "failed to execute command"))
		jobResult <- result
		return
	}

	result.SetData(&remoteResult)
	jobResult <- result
}
