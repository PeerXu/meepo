package config

type Identity struct {
	NoFile bool   `yaml:"no_file" mapstructure:"no_file"`
	File   string `yaml:"file"`
}
