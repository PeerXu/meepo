package config

type API struct {
	Name string `yaml:"name"`

	// http
	Host string `yaml:"host"`
}
