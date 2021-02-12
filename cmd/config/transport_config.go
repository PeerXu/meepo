package config

type TransportConfig struct {
	Name string `yaml:"name"`
}

type ORTCTransportConfig struct {
	Name       string           `yaml:"name"`
	ICEServers []string         `yaml:"iceServers"`
	Signaling  *SignalingConfig `yaml:"signaling"`
	SignalingI interface{}      `yaml:"-"`
}

func (otc *ORTCTransportConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error

	var fotc struct {
		Name       string           `yaml:"name"`
		ICEServers []string         `yaml:"iceServers"`
		Signaling  *SignalingConfig `yaml:"signaling"`
	}

	if err = unmarshal(&fotc); err != nil {
		return err
	}

	switch fotc.Signaling.Name {
	case "redis":
		var frsc struct {
			Signaling RedisSignalingConfig `yaml:"signaling"`
		}
		if err = unmarshal(&frsc); err != nil {
			return err
		}
		otc.SignalingI = &frsc.Signaling
	default:
		return UnsupportedSignalingNameError(fotc.Signaling.Name)
	}

	otc.Name = fotc.Name
	otc.ICEServers = fotc.ICEServers
	otc.Signaling = fotc.Signaling

	return nil
}

func (otc *ORTCTransportConfig) MarshalYAML() (interface{}, error) {
	_otc := struct {
		Name       string      `yaml:"name"`
		ICEServers []string    `yaml:"iceServers"`
		Signaling  interface{} `yaml:"signaling"`
	}{
		Name:       otc.Name,
		ICEServers: otc.ICEServers,
		Signaling:  otc.SignalingI,
	}

	return &_otc, nil
}
