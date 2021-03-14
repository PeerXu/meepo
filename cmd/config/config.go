package config

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

func GetDefaultConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/meepo/meepo.yaml"
	default:
		return "~/.meepo/config.yaml"
	}
}

type Config struct {
	Meepo *MeepoConfig `yaml:"meepo"`
}

func (c *Config) Get(key string) (string, error) {
	switch key {
	case "id":
		return c.Meepo.ID, nil
	case "log.level":
		return c.Meepo.Log.Level, nil
	case "signaling.url":
		return c.Meepo.SignalingI.(*RedisSignalingConfig).URL, nil
	case "transport.iceServers":
		return strings.Join(c.Meepo.TransportI.(*WebrtcTransportConfig).ICEServers, ","), nil
	case "asSignaling":
		return cast.ToString(c.Meepo.AsSignaling), nil
	default:
		return "", UnsupportedGetConfigKeyError(key)
	}
}

func (c *Config) Set(key, val string) error {
	switch key {
	case "id":
		c.Meepo.ID = val
	case "log.level":
		c.Meepo.Log.Level = val
	case "signaling.url":
		c.Meepo.SignalingI.(*RedisSignalingConfig).URL = val
	case "asSignaling":
		c.Meepo.AsSignaling = cast.ToBool(val)
	default:
		return UnsupportedSetConfigKeyError(key)
	}

	return nil
}

func (c *Config) Dump(p string) error {
	var err error

	if p, err = homedir.Expand(p); err != nil {
		return err
	}

	if err = os.MkdirAll(path.Dir(p), 0755); err != nil {
		return err
	}

	buf, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(p, buf, 0644); err != nil {
		return err
	}

	return nil
}

func Load(p string) (config *Config, loaded bool, err error) {
	if p, err = homedir.Expand(p); err != nil {
		return nil, false, err
	}

	config = NewDefaultConfig()

	buf, err := ioutil.ReadFile(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, false, err
		}
		return config, false, nil
	}

	if err = yaml.Unmarshal(buf, config); err != nil {
		return nil, false, err
	}

	return config, true, nil
}

func NewDefaultConfig() *Config {
	return &Config{
		Meepo: &MeepoConfig{
			ID:          "",
			Daemon:      true,
			AsSignaling: false,
			Log: &LogConfig{
				Level: "error",
			},
			Transport: &TransportConfig{
				Name: "webrtc",
			},
			TransportI: &WebrtcTransportConfig{
				Name: "webrtc",
				ICEServers: []string{
					"stun:stun.xten.com:3478",
					"stun:stun.voipbuster.com:3478",
					"stun:stun.sipgate.net:3478",
					"stun:stun.ekiga.net:3478",
					"stun:stun.ideasip.com:3478",
					"stun:stun.schlund.de:3478",
					"stun:stun.voiparound.com:3478",
					"stun:stun.voipbuster.com:3478",
					"stun:stun.voipstunt.com:3478",
					"stun:stun.counterpath.com:3478",
					"stun:stun.1und1.de:3478",
					"stun:stun.gmx.net:3478",
					"stun:stun.callwithus.com:3478",
					"stun:stun.counterpath.net:3478",
					"stun:stun.internetcalls.com:3478",
					"stun:numb.viagenie.ca:3478",
				},
			},
			Signaling: &SignalingConfig{
				Name: "redis",
			},
			SignalingI: &RedisSignalingConfig{
				Name: "redis",
				URL:  "redis://meepo.redis.signaling.peerstud.io:6379/1",
			},
			Api: &ApiConfig{
				Name: "http",
			},
			ApiI: &HttpApiConfig{
				Name: "http",
				Host: "127.0.0.1",
				Port: 12345,
			},
		},
	}
}
