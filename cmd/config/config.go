package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Meepo *MeepoConfig `yaml:"meepo"`
}

func (c *Config) Get(key string) (string, error) {
	switch key {
	case "id":
		return c.Meepo.ID, nil
	case "log.level":
		return c.Meepo.Log.Level, nil
	case "transport.signaling.url":
		return c.Meepo.TransportI.(*ORTCTransportConfig).SignalingI.(*RedisSignalingConfig).URL, nil
	case "transport.iceServers":
		return strings.Join(c.Meepo.TransportI.(*ORTCTransportConfig).ICEServers, ","), nil
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
	case "transport.signaling.url":
		c.Meepo.TransportI.(*ORTCTransportConfig).SignalingI.(*RedisSignalingConfig).URL = val
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
			ID:     "",
			Daemon: true,
			Log: &LogConfig{
				Level: "error",
			},
			Transport: &TransportConfig{
				Name: "ortc",
			},
			TransportI: &ORTCTransportConfig{
				Name: "ortc",
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
				Signaling: &SignalingConfig{
					Name: "redis",
				},
				SignalingI: &RedisSignalingConfig{
					Name: "redis",
					URL:  "redis://meepo.redis.signaling.peerstud.io:6379/1",
				},
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
