package config

type AuthConfig struct {
	Name string `yaml:"name"`
}

type DummyAuthConfig struct {
	Name string `yaml:"name"`
}

type SecretAuthConfig struct {
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
}

func init() {
	RegisterUnmarshalConfigFunc("meepo.auth", "dummy", func(u func(interface{}) error) (interface{}, error) {
		var t struct{ Auth *DummyAuthConfig }
		if err := u(&t); err != nil {
			return nil, err
		}
		return t.Auth, nil
	})
	RegisterUnmarshalConfigFunc("meepo.auth", "secret", func(u func(interface{}) error) (interface{}, error) {
		var t struct{ Auth *SecretAuthConfig }
		if err := u(&t); err != nil {
			return nil, err
		}
		return t.Auth, nil
	})
}
