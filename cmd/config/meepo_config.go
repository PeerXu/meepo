package config

type MeepoConfig struct {
	ID         string           `yaml:"id"`
	Daemon     bool             `yaml:"daemon"`
	Log        *LogConfig       `yaml:"log"`
	Transport  *TransportConfig `yaml:"transport"`
	TransportI interface{}      `yaml:"-"`
	Api        *ApiConfig       `yaml:"api"`
	ApiI       interface{}      `yaml:"-"`
}

func (mc *MeepoConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error

	var fmc struct {
		ID        string           `yaml:"id"`
		Daemon    bool             `yaml:"daemon"`
		Log       *LogConfig       `yaml:"log"`
		Transport *TransportConfig `yaml:"transport"`
		Api       *ApiConfig       `yaml:"api"`
	}

	if err = unmarshal(&fmc); err != nil {
		return err
	}

	switch fmc.Transport.Name {
	case "ortc":
		var fortc struct {
			Transport ORTCTransportConfig `yaml:"transport"`
		}
		if err = unmarshal(&fortc); err != nil {
			return err
		}
		mc.TransportI = &fortc.Transport
	default:
		return UnsupportedTransportNameError(fmc.Transport.Name)
	}

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
	mc.Log = fmc.Log
	mc.Transport = fmc.Transport

	return nil
}

func (mc *MeepoConfig) MarshalYAML() (interface{}, error) {
	_mc := struct {
		ID        string      `yaml:"id"`
		Daemon    bool        `yaml:"daemon"`
		Log       *LogConfig  `yaml:"log"`
		Transport interface{} `yaml:"transport"`
		Api       interface{} `yaml:"api"`
	}{
		ID:        mc.ID,
		Daemon:    mc.Daemon,
		Log:       mc.Log,
		Transport: mc.TransportI,
		Api:       mc.ApiI,
	}

	return &_mc, nil
}
