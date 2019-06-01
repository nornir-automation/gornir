package gornir

type Host struct {
	Hostname string `yaml:"hostname"`
	Port     uint8 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Platform string `yaml:"platform"`
}

type Inventory struct {
	Hosts map[string]*Host
}
