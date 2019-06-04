package gornir_test

import (
	"context"
	"github.com/nornir-automation/gornir/pkg/gornir"
	inv "github.com/nornir-automation/gornir/pkg/plugins/inventory"
	log "github.com/nornir-automation/gornir/pkg/plugins/logger"
	"testing"
)

var (
	file      = "../../examples/hosts.yaml"
	logger    = log.NewLogrus(false)
	noFileErr = "could not read inventory from plugin: problem reading hosts file: open : no such file or directory"
)

func TestRead(t *testing.T) {
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
			plugin := inv.FromYAML{HostsFile: tc.input}

			_, err := gornir.Build(
				gornir.WithInventory(plugin),
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

func TestBuild(t *testing.T) {
	f1 := func(ctx context.Context, h *gornir.Host) bool {
		return h.Hostname == "dev1.group_1" || h.Hostname == "dev4.group_2"
	}
	f2 := func(ctx context.Context, h *gornir.Host) bool {
		return h.Hostname == "uknownk"
	}
	tt := []struct {
		name   string
		input  string
		err    string
		filter gornir.FilterFunc
		length int
	}{
		{name: "With Filter 1", input: file, filter: f1, length: 2},
		{name: "With Filter 2", input: file, filter: f2, length: 0},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			plugin := inv.FromYAML{HostsFile: tc.input}

			original, err := gornir.Build(
				gornir.WithInventory(plugin),
				gornir.WithLogger(logger),
			)
			if err != nil {
				t.Fatalf("could not build a Gornir from file '%s' in Test Case '%s'. Error: '%v'",
					tc.input, tc.name, err)
			}
			filtered, err := original.Build(gornir.WithFilter(tc.filter))
			if err != nil {
				t.Fatalf("could not build a Filtered Gornir in Test Case '%s'. Error: '%v'",
					tc.name, err)
			}
			if len(filtered.Inventory.Hosts) != tc.length {
				t.Fatalf("Inventory Length in Test Case '%s' is %v, want %v",
					tc.name, len(filtered.Inventory.Hosts), tc.length)
			}
		})
	}
}
