package gornir_test

import (
	"context"
	"sync"
	"testing"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"

	"github.com/google/go-cmp/cmp"
)

type dummyTask struct {
}

func (t *dummyTask) Metadata() *gornir.TaskMetadata {
	return nil
}

type dummyTaskResult struct {
}

func (t *dummyTask) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	return dummyTaskResult{}, nil
}

type dummyProcessor struct {
	mux  *sync.Mutex
	data map[string]interface{}
}

func dummy(data map[string]interface{}) *dummyProcessor {
	return &dummyProcessor{
		mux:  &sync.Mutex{},
		data: data,
	}
}

func (r *dummyProcessor) TaskStarted(ctx context.Context, logger gornir.Logger, task gornir.Task) error {
	r.data["started"] = struct{}{}
	return nil
}

func (r *dummyProcessor) TaskCompleted(ctx context.Context, logger gornir.Logger, task gornir.Task) error {
	r.data["completed"] = struct{}{}
	return nil
}

func (r *dummyProcessor) TaskInstanceStarted(ctx context.Context, logger gornir.Logger, host *gornir.Host, task gornir.Task) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	n := host.Hostname + "_started"
	r.data[n] = struct{}{}

	return nil
}

func (r *dummyProcessor) TaskInstanceCompleted(ctx context.Context, logger gornir.Logger, jobResult *gornir.JobResult, host *gornir.Host, task gornir.Task) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	n := host.Hostname + "_completed"
	r.data[n] = struct{}{}

	return nil
}

func TestNoProcessor(t *testing.T) {
	inv := gornir.Inventory{
		Hosts: map[string]*gornir.Host{
			"host1": {Hostname: "host1"},
			"host2": {Hostname: "host2"},
		},
	}
	log := logger.NewLogrus(false)
	rnr := runner.Sorted()
	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	_, err := gr.RunSync(context.Background(), &dummyTask{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestOneProcessor(t *testing.T) {
	inv := gornir.Inventory{
		Hosts: map[string]*gornir.Host{
			"host1": {Hostname: "host1"},
			"host2": {Hostname: "host2"},
		},
	}
	log := logger.NewLogrus(false)
	rnr := runner.Sorted()
	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	data := make(map[string]interface{})

	expected := map[string]interface{}{
		"completed":       struct{}{},
		"host1_completed": struct{}{},
		"host1_started":   struct{}{},
		"host2_completed": struct{}{},
		"host2_started":   struct{}{},
		"started":         struct{}{},
	}

	_, err := gr.WithProcessors(gornir.Processors{dummy(data)}).RunSync(context.Background(), &dummyTask{})
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(data, expected) {
		t.Error(cmp.Diff(data, expected))
	}
}

func TestMultipleProcessor(t *testing.T) {
	inv := gornir.Inventory{
		Hosts: map[string]*gornir.Host{
			"host1": {Hostname: "host1"},
			"host2": {Hostname: "host2"},
		},
	}
	log := logger.NewLogrus(false)
	rnr := runner.Sorted()
	gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

	data1 := make(map[string]interface{})
	data2 := make(map[string]interface{})

	expected := map[string]interface{}{
		"completed":       struct{}{},
		"host1_completed": struct{}{},
		"host1_started":   struct{}{},
		"host2_completed": struct{}{},
		"host2_started":   struct{}{},
		"started":         struct{}{},
	}

	_, err := gr.WithProcessors(gornir.Processors{dummy(data1), dummy(data2)}).RunSync(context.Background(), &dummyTask{})
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(data1, expected) {
		t.Error(cmp.Diff(data2, expected))
	}

	if !cmp.Equal(data2, expected) {
		t.Error(cmp.Diff(data2, expected))
	}
}
