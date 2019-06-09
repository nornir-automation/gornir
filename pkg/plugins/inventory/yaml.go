package inventory

import (
	"io/ioutil"

	"github.com/nornir-automation/gornir/pkg/gornir"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// FromYAML satisfies the InventoryPlugin interface for YAML files.
type FromYAML struct {
	HostsFile string
}

// Create parses the content of a YAML file follwoing the same structure
// as the structs, but in lower case to create an Inventory. For instance:
//     dev1.group_1:
//         port: 22
//         hostname: dev1.group_1
//         username: root
//         password: docker
//
//     dev2.group_1:
//         port: 22
//         hostname: dev2.group_1
//         username: root
//         password: docker
func (f FromYAML) Create() (gornir.Inventory, error) {
	b, err := ioutil.ReadFile(f.HostsFile)
	if err != nil {
		return gornir.Inventory{}, errors.Wrap(err, "problem reading hosts file")
	}
	hosts := make(map[string]*gornir.Host)
	err = yaml.Unmarshal(b, hosts)
	if err != nil {
		return gornir.Inventory{}, errors.Wrap(err, "problem unmarshalling yaml")
	}

	return gornir.Inventory{
		Hosts: hosts,
	}, nil
}
