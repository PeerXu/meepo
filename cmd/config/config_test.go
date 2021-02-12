package config_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/PeerXu/meepo/cmd/config"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestEncodeDecodeORTCTransportConfig() {
	in := &config.ORTCTransportConfig{
		Name: "ortc",
		Signaling: &config.SignalingConfig{
			Name: "redis",
		},
		SignalingI: &config.RedisSignalingConfig{
			Name: "redis",
			URL:  "redis://localhost:6379/1",
		},
	}

	buf, err := yaml.Marshal(in)
	s.Require().Nil(err)

	var out config.ORTCTransportConfig
	err = yaml.Unmarshal(buf, &out)
	s.Require().Nil(err)

	s.Equal("redis", out.Signaling.Name)
	s.Equal("ortc", out.Name)

	si, ok := out.SignalingI.(*config.RedisSignalingConfig)
	s.Require().True(ok)
	s.Equal("redis", si.Name)
	s.Equal("redis://localhost:6379/1", si.URL)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
