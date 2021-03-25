package config

type TransportConfig struct {
	Name string `yaml:"name"`
}

type WebrtcTransportConfig struct {
	Name       string   `yaml:"name"`
	ICEServers []string `yaml:"iceServers"`
}

func init() {
	RegisterUnmarshalConfigFunc("meepo.transport", "webrtc", func(u func(interface{}) error) (interface{}, error) {
		var t struct{ Transport *WebrtcTransportConfig }
		if err := u(&t); err != nil {
			return nil, err
		}
		return t.Transport, nil
	})
}
