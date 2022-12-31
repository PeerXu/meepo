package config

type Tracker struct {
	Name string `yaml:"name"`

	// rpc
	CallerName string `yaml:"callerName"`
	Addr       string `yaml:"addr"`
	Host       string `yaml:"host,omitempty"`
}
