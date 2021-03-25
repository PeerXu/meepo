package config

type SignalingConfig struct {
	Name string `yaml:"name"`
}

type RedisSignalingConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func init() {
	RegisterUnmarshalConfigFunc("meepo.signaling", "redis", func(u func(interface{}) error) (interface{}, error) {
		var t struct{ Signaling *RedisSignalingConfig }
		if err := u(&t); err != nil {
			return nil, err
		}
		return t.Signaling, nil
	})
}
