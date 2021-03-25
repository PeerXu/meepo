package config

type ApiConfig struct {
	Name string `yaml:"name"`
}

type HttpApiConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
}

func init() {
	RegisterUnmarshalConfigFunc("meepo.api", "http", func(u func(interface{}) error) (interface{}, error) {
		var t struct {
			Api *HttpApiConfig
		}

		if err := u(&t); err != nil {
			return nil, err
		}

		return t.Api, nil
	})
}
