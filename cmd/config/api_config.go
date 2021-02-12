package config

type ApiConfig struct {
	Name string `yaml:"name"`
}

type HttpApiConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
}
