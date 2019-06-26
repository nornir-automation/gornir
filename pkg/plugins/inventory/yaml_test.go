package inventory_test

import (
	"testing"

	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
)

var (
	file          = "testdata/hosts.yaml"
	noFileErr     = "problem reading hosts file: open : no such file or directory"
	incorrectYAML = "testdata/inchosts"
	incYAMLErr    = "problem unmarshalling yaml: yaml: line 9: did not find expected key"
)

func TestCreate(t *testing.T) {
	tt := []struct {
		name  string
		input string
		err   string
	}{
		{name: "From YAML file", input: file},
		{name: "From no file", input: "", err: noFileErr},
		{name: "From Incorrect YAML", input: incorrectYAML, err: incYAMLErr},
	}
	for _, tc := range tt {
		tc := tc // lock the variable. This is a problem of golint, we don't need this here.
		t.Run(tc.name, func(t *testing.T) {
			plugin := inventory.FromYAML{HostsFile: tc.input}
			_, err := plugin.Create()

			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("could not read an inventory from file '%s' in Test Case '%s'. Error: '%v'",
						tc.input, tc.name, err)
				}
			}
		})
	}
}

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		plugin := inventory.FromYAML{HostsFile: file}
		_, err := plugin.Create()
		if err != nil {
			b.Fatalf("could not read an inventory from file '%s' in Benchmark", file)
		}
		// _ = gornir.New().WithInventory(inv)
	}
}
