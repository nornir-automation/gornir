package gornir_test

import (
	"context"
	"testing"
	"time"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
)

var (
	file      = "../../examples/hosts.yaml"
	log       = logger.NewLogrus(false)
	noFileErr = "problem reading hosts file: open : no such file or directory"
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
		tc := tc // lock the variable
		t.Run(tc.name, func(t *testing.T) {
			plugin := inventory.FromYAML{HostsFile: tc.input}
			inv, err := plugin.Create()

			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("could not read an inventory from file '%s' in Test Case '%s'. Error: '%v'",
						tc.input, tc.name, err)
				}
			}

			_ = gornir.New().WithInventory(inv).WithLogger(log)

		})
	}
}

func TestBuild(t *testing.T) {
	f1 := func(h *gornir.Host) bool {
		return h.Hostname == "dev1.group_1" || h.Hostname == "dev4.group_2"
	}
	f2 := func(h *gornir.Host) bool {
		return h.Hostname == "uknownk"
	}
	tt := []struct {
		name   string
		input  string
		filter gornir.FilterFunc
		length int
	}{
		{name: "With Filter 1", input: file, filter: f1, length: 2},
		{name: "With Filter 2", input: file, filter: f2, length: 0},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			plugin := inventory.FromYAML{HostsFile: tc.input}
			inv, err := plugin.Create()

			original := gornir.New().WithInventory(inv).WithLogger(log)
			olen := len(original.Inventory.Hosts)

			if err != nil {
				t.Fatalf("could not build a Gornir from file '%s' in Test Case '%s'. Error: '%v'",
					tc.input, tc.name, err)
			}
			filtered := original.Filter(tc.filter)
			if err != nil {
				t.Fatalf("could not build a Filtered Gornir in Test Case '%s'. Error: '%v'",
					tc.name, err)
			}
			if len(filtered.Inventory.Hosts) != tc.length {
				t.Fatalf("Filtered Inventory Length in Test Case '%s' is %v, want %v",
					tc.name, len(filtered.Inventory.Hosts), tc.length)
			}
			if len(original.Inventory.Hosts) != olen {
				t.Fatalf("Oringinal Inventory Length in Test Case '%s' is %v, want %v",
					tc.name, len(original.Inventory.Hosts), olen)
			}
		})
	}
}

type slowTask struct {
}

func (t *slowTask) Metadata() *gornir.TaskMetadata {
	return nil
}

type slowTaskResult struct {
}

func (t *slowTask) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	select {
	case <-time.After(1 * time.Second):
		return slowTaskResult{}, nil
	case <-ctx.Done():
		return slowTaskResult{}, ctx.Err()
	}
}

func TestContextCancelSync(t *testing.T) {
	inv := gornir.Inventory{
		Hosts: map[string]*gornir.Host{
			"host1": {},
		},
	}
	log := logger.NewLogrus(false)
	rnr := runner.Sorted()
	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	res, err := gr.RunSync(ctx, &slowTask{})
	if err != nil {
		t.Error(err)
	}
	r := <-res
	if r.Err() != context.DeadlineExceeded {
		t.Errorf("error should be 'context deadline exceeded'. Got: %s", r.Err())
	}
}

func TestContextCancelASync(t *testing.T) {
	inv := gornir.Inventory{
		Hosts: map[string]*gornir.Host{
			"host1": {},
		},
	}
	log := logger.NewLogrus(false)
	rnr := runner.Sorted()
	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	res := make(chan *gornir.JobResult, len(gr.Inventory.Hosts))
	defer close(res)

	err := gr.RunAsync(ctx, &slowTask{}, res)
	if err != nil {
		t.Error(err)
	}

	err = rnr.Wait()
	if err != nil {
		t.Error(err)
	}

	r := <-res
	if r.Err() != context.DeadlineExceeded {
		t.Errorf("error should be 'context deadline exceeded'. Got: %s", r.Err())
	}
}
