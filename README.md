[![GoDoc](https://godoc.org/github.com/nornir-automation/gornir?status.svg)](http://godoc.org/github.com/nornir-automation/gornir)
[![Build Status](https://travis-ci.org/nornir-automation/gornir.svg?branch=master)](https://travis-ci.org/nornir-automation/gornir)
[![codecov](https://codecov.io/gh/nornir-automation/gornir/branch/master/graph/badge.svg)](https://codecov.io/gh/nornir-automation/gornir)
[![Go Report Card](https://goreportcard.com/badge/github.com/nornir-automation/gornir)](https://goreportcard.com/report/github.com/nornir-automation/gornir)

gornir
======

Gornir is a pluggable framework with inventory management to help operate collections of devices. It's similar to [nornir](https://github.com/nornir-automation/nornir/) but in golang.

The goal is to be able to operate on many devices with little effort. For instance:

```go
package main

import (
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	logger := logger.NewLogrus(false)

	inventory, err := inventory.FromYAMLFile("/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	results, err := gr.RunSync(
		"What's my ip?",
		runner.Parallel(),
		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
	)
	if err != nil {
		logger.Fatal(err)
	}
	output.RenderResults(os.Stdout, results, true)
}
```

would render:

```bash
# What's my ip?
@ dev5.no_group
  - err: failed to dial: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain

@ dev1.group_1
 * Stdout: 10.21.33.101/24

 * Stderr:
  - err: <nil>

@ dev2.group_1
 * Stdout: 10.21.33.102/24

 * Stderr:
  - err: <nil>

@ dev3.group_2
 * Stdout: 10.21.33.103/24

 * Stderr:
  - err: <nil>

@ dev4.group_2
 * Stdout: 10.21.33.104/24

 * Stderr:
  - err: <nil>
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
