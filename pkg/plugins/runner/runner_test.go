package runner

import (
	"context"
	"time"

	"github.com/nornir-automation/gornir/pkg/gornir"
)

var (
	testHosts = map[string]*gornir.Host{
		"dev1": {Hostname: "dev1"},
		"dev2": {Hostname: "dev2"},
		"dev3": {Hostname: "dev3"},
		"dev4": {Hostname: "dev4"},
	}
)

type testTaskSleep struct {
	sleepDuration time.Duration
}

type testTaskSleepResults struct {
	success bool
}

func (t *testTaskSleep) Run(ctx context.Context, host *gornir.Host) (interface{}, error) {
	time.Sleep(t.sleepDuration)
	return testTaskSleepResults{success: true}, nil
}
