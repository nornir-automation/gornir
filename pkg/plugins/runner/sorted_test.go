package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
)

// TestSorted runs test func and verifies the hosts are executed
// in alphabetical order by checking the results come in the right
// order
func TestSorted(t *testing.T) {
	testCases := []struct {
		name          string
		expected      []string
		sleepDuration time.Duration
	}{
		{
			name:          "simple test",
			expected:      []string{"dev1", "dev2", "dev3", "dev4"},
			sleepDuration: 1 * time.Millisecond,
		},
	}

	testHosts = map[string]*gornir.Host{
		"dev1": {Hostname: "dev1"},
		"dev2": {Hostname: "dev2"},
		"dev3": {Hostname: "dev3"},
		"dev4": {Hostname: "dev4"},
	}

	for _, tc := range testCases {
		tc := tc
		results := make(chan *gornir.JobResult, len(testHosts))
		t.Run(tc.name, func(t *testing.T) {
			rnr := runner.Sorted()
			if err := rnr.Run(
				context.Background(),
				&testTaskSleep{sleepDuration: tc.sleepDuration},
				testHosts,
				gornir.NewJobParameters("test", logger.NewLogrus(false)),
				results,
			); err != nil {
				t.Fatal(err)
			}
			if err := rnr.Wait(); err != nil {
				t.Fatal(err)
			}
			close(results)

			// let's process the results and turn it into a list so we can
			// compare with our expected value
			got := make([]string, len(testHosts))
			i := 0
			for res := range results {
				got[i] = res.JobParameters().Host().Hostname
				i++
			}
			if !cmp.Equal(got, tc.expected) {
				t.Error(cmp.Diff(got, tc.expected))
			}
		})
	}

}
