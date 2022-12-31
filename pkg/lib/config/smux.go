package config

type Smux struct {
	Disable          bool `yaml:"disable,omitempty"`
	Version          int  `yaml:"version,omitempty"`
	BufferSize       int  `yaml:"bufferSize,omitempty"`
	StreamBufferSize int  `yaml:"streamBufferSize,omitempty"`
	Nocomp           bool `yaml:"nocomp,omitempty"`
}
