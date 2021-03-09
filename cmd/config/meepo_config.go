package config

type MeepoConfig struct {
	ID          string           `yaml:"id"`
	Daemon      bool             `yaml:"daemon"`
	AsSignaling bool             `yaml:"asSignaling"`
	Log         *LogConfig       `yaml:"log"`
	Transport   *TransportConfig `yaml:"transport"`
	TransportI  interface{}      `yaml:"-"`
	Signaling   *SignalingConfig `yaml:"signaling"`
	SignalingI  interface{}      `yaml:"-"`
	Api         *ApiConfig       `yaml:"api"`
	ApiI        interface{}      `yaml:"-"`
}

func (mc *MeepoConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error

	var fmc struct {
		ID          string           `yaml:"id"`
		Daemon      bool             `yaml:"daemon"`
		AsSignaling bool             `yaml:"asSignaling"`
		Log         *LogConfig       `yaml:"log"`
		Transport   *TransportConfig `yaml:"transport"`
		Signaling   *SignalingConfig `yaml:"signaling"`
		Api         *ApiConfig       `yaml:"api"`
	}

	if err = unmarshal(&fmc); err != nil {
		return err
	}

	mc.Transport = fmc.Transport
	switch fmc.Transport.Name {
	case "webrtc":
		var fwrtc struct {
			Transport WebrtcTransportConfig `yaml:"transport"`
		}
		if err = unmarshal(&fwrtc); err != nil {
			return err
		}
		mc.TransportI = &fwrtc.Transport
	default:
		return UnsupportedTransportNameError(fmc.Transport.Name)
	}

	mc.Signaling = fmc.Signaling
	switch fmc.Signaling.Name {
	case "redis":
		var frsc struct {
			Signaling RedisSignalingConfig `yaml:"signaling"`
		}
		if err = unmarshal(&frsc); err != nil {
			return err
		}
		mc.SignalingI = &frsc.Signaling
	}

	mc.Api = fmc.Api
	switch fmc.Api.Name {
	case "http":
		var fhac struct {
			Api HttpApiConfig `yaml:"api"`
		}
		if err = unmarshal(&fhac); err != nil {
			return err
		}
		mc.ApiI = &fhac.Api
	default:
		return UnsupportedApiNameError(fmc.Api.Name)
	}

	mc.ID = fmc.ID
	mc.Daemon = fmc.Daemon
	mc.AsSignaling = fmc.AsSignaling
	mc.Log = fmc.Log

	return nil
}

func (mc *MeepoConfig) MarshalYAML() (interface{}, error) {
	_mc := struct {
		ID          string      `yaml:"id"`
		Daemon      bool        `yaml:"daemon"`
		AsSignaling bool        `yaml:"asSignaling"`
		Log         *LogConfig  `yaml:"log"`
		Transport   interface{} `yaml:"transport"`
		Signaling   interface{} `yaml:"signaling"`
		Api         interface{} `yaml:"api"`
	}{
		ID:          mc.ID,
		Daemon:      mc.Daemon,
		AsSignaling: mc.AsSignaling,
		Log:         mc.Log,
		Transport:   mc.TransportI,
		Signaling:   mc.SignalingI,
		Api:         mc.ApiI,
	}

	return &_mc, nil
}
