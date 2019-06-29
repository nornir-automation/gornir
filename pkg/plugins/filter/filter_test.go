package filter

import (
	"errors"
	"sort"
	"testing"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/google/go-cmp/cmp"
)

var (
	err1 = errors.New("an error")
	err2 = errors.New("another error")
)

func TestFilters(t *testing.T) {
	tt := []struct {
		name     string
		filter   gornir.FilterFunc
		expected []string
	}{
		{
			"Successful",
			Not(Errored),
			[]string{"dev1", "dev3"},
		},
		{
			"Errored",
			Errored,
			[]string{"dev2", "dev4"},
		},
		{
			"WithError",
			WithError(err1),
			[]string{"dev2"},
		},
		{
			"WithoutError",
			Not(WithError(err1)),
			[]string{"dev1", "dev3", "dev4"},
		},
		{
			"And_Empty",
			And(),
			[]string{},
		},
		{
			"And_Pass",
			And(WithHostname("dev1"), Not(Errored)),
			[]string{"dev1"},
		},
		{
			"And_Failed",
			And(WithHostname("dev2"), Not(Errored)),
			[]string{},
		},
		{
			"Or_Empty",
			Or(),
			[]string{},
		},
		{
			"Or_Pass",
			Or(WithHostname("dev1"), Errored),
			[]string{"dev1", "dev2", "dev4"},
		},
		{
			"Or_Failed",
			Or(WithHostname("dev5"), WithError(errors.New("yet another oerr"))),
			[]string{},
		},
	}

	inv := &gornir.Inventory{
		Hosts: map[string]*gornir.Host{
			"dev1": {Hostname: "dev1"},
			"dev2": {Hostname: "dev2"},
			"dev3": {Hostname: "dev3"},
			"dev4": {Hostname: "dev4"},
		},
	}
	inv.Hosts["dev2"].SetErr(err1)
	inv.Hosts["dev4"].SetErr(err2)

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotInv := inv.Filter(tc.filter)

			got := make([]string, len(gotInv.Hosts))
			i := 0
			for h := range gotInv.Hosts {
				got[i] = h
				i++
			}
			sort.Strings(got)
			if !cmp.Equal(got, tc.expected) {
				t.Error(cmp.Diff(got, tc.expected))
			}
		})
	}
}
