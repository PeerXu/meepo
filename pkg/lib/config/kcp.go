package config

type Kcp struct {
	Disable     bool   `yaml:"disable,omitempty"`
	Preset      string `yaml:"preset,omitempty"`
	Crypt       string `yaml:"crypt,omitempty"`
	Key         string `yaml:"key,omitempty"`
	Mtu         int    `yaml:"mtu,omitempty"`
	Sndwnd      int    `yaml:"sndwnd,omitempty"`
	Rcvwnd      int    `yaml:"rcvwnd,omitempty"`
	DataShard   int    `yaml:"dataShard,omitempty"`
	ParityShard int    `yaml:"parityShard,omitempty"`
}
