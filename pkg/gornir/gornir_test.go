package gornir_test

import (
	"github.com/nornir-automation/gornir/pkg/gornir"
	log "github.com/nornir-automation/gornir/pkg/plugins/logger"
	"testing"
)

var (
	file      = "../../examples/hosts.yaml"
	logger    = log.NewLogrus(false)
	noFileErr = "could not read inventory from file : problem reading hosts file: open : no such file or directory"
)

func TestBuild(t *testing.T) {
	tt := []struct {
		name  string
		input string
		err   string
	}{
		{name: "From YAML", input: file},
		{name: "From no file", input: "", err: noFileErr},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Instantiate Gornir
			_, err := gornir.Build(
				gornir.WithInventory(tc.input),
				gornir.WithLogger(logger),
			)
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("could not build a Gornir from file '%s' in Test Case '%s'. Error: '%v'",
						tc.input, tc.name, err)
				}
			}
		})
	}
}
