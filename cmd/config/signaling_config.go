package config

type SignalingConfig struct {
	Name string `yaml:"name"`
}

type RedisSignalingConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}
