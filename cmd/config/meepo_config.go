package config

type MeepoConfig struct {
	ID          string           `yaml:"id"`
	Daemon      bool             `yaml:"daemon"`
	AsSignaling bool             `yaml:"asSignaling"`
	Log         *LogConfig       `yaml:"log"`
	Proxy       *ProxyConfig     `yaml:"proxy"`
	Auth        *AuthConfig      `yaml:"auth"`
	AuthI       interface{}      `yaml:"-"`
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
		Proxy       *ProxyConfig     `yaml:"proxy"`
		Auth        *AuthConfig      `yaml:"auth"`
		Transport   *TransportConfig `yaml:"transport"`
		Signaling   *SignalingConfig `yaml:"signaling"`
		Api         *ApiConfig       `yaml:"api"`
	}

	if err = unmarshal(&fmc); err != nil {
		return err
	}

	mc.Auth = fmc.Auth
	if fmc.Auth != nil {
		if mc.AuthI, err = unmarshalConfig("meepo.auth", fmc.Auth.Name, unmarshal); err != nil {
			return err
		}
	}

	mc.Transport = fmc.Transport
	if fmc.Transport != nil {
		if mc.TransportI, err = unmarshalConfig("meepo.transport", fmc.Transport.Name, unmarshal); err != nil {
			return err
		}
	}

	mc.Signaling = fmc.Signaling
	if fmc.Signaling != nil {
		if mc.SignalingI, err = unmarshalConfig("meepo.signaling", fmc.Signaling.Name, unmarshal); err != nil {
			return err
		}
	}

	mc.Api = fmc.Api
	if fmc.Api != nil {
		if mc.ApiI, err = unmarshalConfig("meepo.api", fmc.Api.Name, unmarshal); err != nil {
			return err
		}
	}

	mc.ID = fmc.ID
	mc.Daemon = fmc.Daemon
	mc.AsSignaling = fmc.AsSignaling
	mc.Log = fmc.Log
	mc.Proxy = fmc.Proxy

	return nil
}

func (mc *MeepoConfig) MarshalYAML() (interface{}, error) {
	_mc := struct {
		ID          string       `yaml:"id"`
		Daemon      bool         `yaml:"daemon"`
		AsSignaling bool         `yaml:"asSignaling"`
		Log         *LogConfig   `yaml:"log"`
		Proxy       *ProxyConfig `yaml:"proxy"`
		Auth        interface{}  `yaml:"auth"`
		Transport   interface{}  `yaml:"transport"`
		Signaling   interface{}  `yaml:"signaling"`
		Api         interface{}  `yaml:"api"`
	}{
		ID:          mc.ID,
		Daemon:      mc.Daemon,
		AsSignaling: mc.AsSignaling,
		Log:         mc.Log,
		Proxy:       mc.Proxy,
		Auth:        mc.AuthI,
		Transport:   mc.TransportI,
		Signaling:   mc.SignalingI,
		Api:         mc.ApiI,
	}

	return &_mc, nil
}
