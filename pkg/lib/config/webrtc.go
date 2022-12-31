package config

type Webrtc struct {
	IceServers     []string `yaml:"iceServers,omitempty"`
	RecvBufferSize uint32   `yaml:"recvBufferSize,omitempty"`
}
