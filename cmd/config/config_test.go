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

func (s *ConfigTestSuite) TestEncodeDecodeConfig() {
	in := config.NewDefaultConfig()

	buf, err := yaml.Marshal(in)
	s.Require().Nil(err)

	var out config.Config
	err = yaml.Unmarshal(buf, &out)
	s.Require().Nil(err)

	{
		ai, ok := out.Meepo.AuthI.(*config.DummyAuthConfig)
		s.Require().True(ok)
		s.Equal("dummy", ai.Name)
	}
	{
		ti, ok := out.Meepo.TransportI.(*config.WebrtcTransportConfig)
		s.Require().True(ok)
		s.Equal("webrtc", ti.Name)
	}
	{
		si, ok := out.Meepo.SignalingI.(*config.RedisSignalingConfig)
		s.Require().True(ok)
		s.Equal("redis", si.Name)
	}
	{
		ai, ok := out.Meepo.ApiI.(*config.HttpApiConfig)
		s.Require().True(ok)
		s.Equal("http", ai.Name)
	}
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
