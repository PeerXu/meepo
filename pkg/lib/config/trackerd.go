package config

type Trackerd struct {
	Name string `yaml:"name"`

	// rpc
	ServerName string `yaml:"serverName,omitempty"`
	Host       string `yaml:"host,omitempty"`
}
