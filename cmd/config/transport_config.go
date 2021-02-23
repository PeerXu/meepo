package config

type TransportConfig struct {
	Name string `yaml:"name"`
}

type WebrtcTransportConfig struct {
	Name       string   `yaml:"name"`
	ICEServers []string `yaml:"iceServers"`
}

func (wtc *WebrtcTransportConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error

	var fwtc struct {
		Name       string   `yaml:"name"`
		ICEServers []string `yaml:"iceServers"`
	}

	if err = unmarshal(&fwtc); err != nil {
		return err
	}

	wtc.Name = fwtc.Name
	wtc.ICEServers = fwtc.ICEServers

	return nil
}

func (wtc *WebrtcTransportConfig) MarshalYAML() (interface{}, error) {
	_wtc := struct {
		Name       string   `yaml:"name"`
		ICEServers []string `yaml:"iceServers"`
	}{
		Name:       wtc.Name,
		ICEServers: wtc.ICEServers,
	}

	return &_wtc, nil
}
