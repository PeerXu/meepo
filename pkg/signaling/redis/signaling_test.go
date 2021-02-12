package redis_signaling_test

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/suite"

	"github.com/PeerXu/meepo/pkg/signaling"
	redis_signaling "github.com/PeerXu/meepo/pkg/signaling/redis"
)

func newEngineOptions(id string) []signaling.NewEngineOption {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)

	return []signaling.NewEngineOption{
		signaling.WithID(id),
		redis_signaling.WithURL(viper.GetString("redis_url")),
		signaling.WithLogger(logger),
	}
}

type RedisEngineTestSuite struct {
	suite.Suite
}

func (s *RedisEngineTestSuite) TestRedisEngine() {
	ea, err := redis_signaling.NewRedisEngine(newEngineOptions("a")...)
	s.Require().Nil(err)
	defer ea.Close()

	ea.OnWire(func(in *signaling.Descriptor) (*signaling.Descriptor, error) {
		ud := objx.New(in.UserData)
		s.Equal(1, cast.ToInt(ud.Get("a").Inter()))

		return in, nil
	})

	time.Sleep(1 * time.Millisecond)
	d, err := ea.Wire(
		&signaling.Descriptor{ID: "a"},
		&signaling.Descriptor{
			ID: "a",
			UserData: map[string]interface{}{
				"a": 1,
			},
		},
	)
	s.Require().Nil(err)
	s.Equal(1, cast.ToInt(objx.New(d.UserData).Get("a").Inter()))
}

func TestRedisEngineTestSuite(t *testing.T) {
	suite.Run(t, new(RedisEngineTestSuite))
}

func init() {
	viper.SetEnvPrefix("mpt")
	viper.BindEnv("redis_url")
	viper.SetDefault("redis_url", "redis://127.0.0.1:6379/0")
}
