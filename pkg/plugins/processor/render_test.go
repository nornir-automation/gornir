package processor_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/processor"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("test.update", false, "update golden files")

func compareWithGolden(t *testing.T, got []byte, goldenPath string) {
	if *update {
		if err := ioutil.WriteFile(goldenPath, got, 0644); err != nil {
			t.Fatal(err)
		}
	}

	expected, err := ioutil.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(got, expected) {
		t.Error(cmp.Diff(string(got), string(expected)))
	}
}

type dummyTask struct {
	meta *gornir.TaskMetadata
}

func (t *dummyTask) Metadata() *gornir.TaskMetadata {
	return t.meta
}

type dummyTaskResult struct {
}

func (r dummyTaskResult) String() string {
	return "  - done!"
}

func (t *dummyTask) Run(ctx context.Context, logger gornir.Logger, host *gornir.Host) (gornir.TaskInstanceResult, error) {
	if host.Hostname == "host2" {
		return dummyTaskResult{}, errors.New("some error")
	}
	return dummyTaskResult{}, nil
}

func TestRender(t *testing.T) {
	cases := []struct {
		name       string
		goldenPath string
		task       gornir.Task
		color      bool
	}{
		{
			name:       "color_task_no_name",
			goldenPath: filepath.Join("testdata", "render", "color_task_no_name.golden"),
			task:       &dummyTask{},
			color:      true,
		},
		{
			name:       "no_color_task_no_name",
			goldenPath: filepath.Join("testdata", "render", "no_color_task_no_name.golden"),
			task:       &dummyTask{},
			color:      false,
		},
		{
			name:       "color_task_with_id",
			goldenPath: filepath.Join("testdata", "render", "color_task_with_id.golden"),
			task: &dummyTask{
				meta: &gornir.TaskMetadata{Identifier: "Some task name"},
			},
			color: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			inv := gornir.Inventory{
				Hosts: map[string]*gornir.Host{
					"host1": {Hostname: "host1"},
					"host2": {Hostname: "host2"},
				},
			}
			log := logger.NewLogrus(false)
			rnr := runner.Sorted()
			gr := gornir.New().WithInventory(inv).WithLogger(log).WithRunner(rnr)

			b := []byte{}
			buf := bytes.NewBuffer(b)

			_, err := gr.WithProcessor(processor.Render(buf, tc.color)).RunSync(tc.task)
			if err != nil {
				t.Fatal(err)
			}

			compareWithGolden(t, buf.Bytes(), tc.goldenPath)

		})
	}
}
