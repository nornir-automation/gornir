package inventory

import (
	"io/ioutil"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"gopkg.in/yaml.v2"
	"github.com/pkg/errors"
)

func FromYAMLFile(hostsFile string) (*gornir.Inventory, error) {
	b, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		return &gornir.Inventory{}, errors.Wrap(err, "problem reading hosts file")
	}
	hosts := make(map[string]*gornir.Host)
	err = yaml.Unmarshal(b, hosts)
	if err != nil {
		return &gornir.Inventory{}, errors.Wrap(err, "problem unmarshalling yaml")
	}

	return &gornir.Inventory{
		Hosts: hosts,
	}, nil
}
