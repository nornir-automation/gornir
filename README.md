[![GoDoc](https://godoc.org/github.com/nornir-automation/gornir?status.svg)](http://godoc.org/github.com/nornir-automation/gornir)
[![Build Status](https://travis-ci.com/nornir-automation/gornir.svg?branch=master)](https://travis-ci.com/nornir-automation/gornir)
[![codecov](https://codecov.io/gh/nornir-automation/gornir/branch/master/graph/badge.svg)](https://codecov.io/gh/nornir-automation/gornir)
[![Go Report Card](https://goreportcard.com/badge/github.com/nornir-automation/gornir)](https://goreportcard.com/report/github.com/nornir-automation/gornir)

gornir
======

Gornir is a pluggable framework with inventory management to help operate collections of devices. It's similar to [nornir](https://github.com/nornir-automation/nornir/) but in golang.

The goal is to be able to operate on many devices with little effort. For instance:

```go
package main

import (
	"context"
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/connection"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	log := logger.NewLogrus(false)

	file := "/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml"
	plugin := inventory.FromYAML{HostsFile: file}
	inv, err := plugin.Create()
	if err != nil {
		log.Fatal(err)
	}

	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(runner.Parallel())

	// Open an SSH connection towards the devices
	results, err := gr.RunSync(
		context.Background(),
		&connection.SSHOpen{},
	)
	if err != nil {
		log.Fatal(err)
	}

	// defer closing the SSH connection we just opened
	defer func() {
		results, err = gr.RunSync(
			context.Background(),
			&connection.SSHClose{},
		)
		if err != nil {
			log.Fatal(err)
		}
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
}
```

would render:

```bash
# What is my ip?
@ dev5.no_group
  - err: failed to retrieve connection: couldn't find connection

@ dev1.group_1
  - stdout: 10.21.33.101/24

  - stderr:
@ dev6.no_group
  - stdout: 10.21.33.106/24

  - stderr:
@ dev4.group_2
  - stdout: 10.21.33.104/24

  - stderr:
@ dev3.group_2
  - stdout: 10.21.33.103/24

  - stderr:
@ dev2.group_1
  - stdout: 10.21.33.102/24

  - stderr:
```

## Examples

You can see more examples in the [examples](examples) folder and run them with [Docker-Compose](https://docs.docker.com/compose/install/) as follows:

1. Create a development enviroment

```bash
make start-dev-env
```

2. Run any of the examples in the [examples](examples) folder with `make example`. Specify the name of the example with `EXAMPLE`; for instance `2_simple_with_filter`.

```bash
make example EXAMPLE=2_simple_with_filter
```

3. After you are done, make sure you stop the development enviroment

```bash
make stop-dev-env
```

The project is still work in progress and feedback/help is welcomed.
