package config

type AclConfig struct {
	Allows []string `yaml:"allows,omitempty"`
	Blocks []string `yaml:"blocks,omitempty"`
}
