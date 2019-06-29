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

// TestParallel is going to check that the func runs on all hosts
// and that it takes less than X time to complete. The test func
// basically will sleep for N ms and given we are using goroutines
// the completion should only be slightly above it even though
// we are sleeping once per device
func TestParallel(t *testing.T) {
	testCases := []struct {
		name          string
		expected      map[string]bool
		sleepDuration time.Duration
	}{
		{
			name:          "simple test",
			expected:      map[string]bool{"dev1": true, "dev2": true, "dev3": true, "dev4": true},
			sleepDuration: 200 * time.Millisecond,
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
			rnr := runner.Parallel()
			startTime := time.Now()
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

			// let's process the results and turn it into a map so we can
			// compare with our expected value
			got := make(map[string]bool)
			for res := range results {
				got[res.JobParameters().Host().Hostname] = res.Data().(*testTaskSleepResults).success
			}
			if !cmp.Equal(got, tc.expected) {
				t.Error(cmp.Diff(got, tc.expected))
			}
			// now we check test took what we expected
			if time.Since(startTime) > (tc.sleepDuration + time.Millisecond*50) {
				t.Errorf("test took to long, parallelization might not be working: %v\n", time.Since(startTime).Seconds())
			}
		})
	}

}
