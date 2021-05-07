package config

type ProxyConfig struct {
	Socks5 *Socks5Config `yaml:"socks5"`
}

type Socks5Config struct {
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
}
