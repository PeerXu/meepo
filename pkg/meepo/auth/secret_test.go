package auth_test

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/suite"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
	mrand "github.com/PeerXu/meepo/pkg/util/random"
)

type SecretAuthEngineTestSuite struct {
	suite.Suite

	e *auth.SecretEngine
}

func (s *SecretAuthEngineTestSuite) SetupTest() {
	mrand.Random.Seed(42)
	e, err := auth.NewEngine("secret",
		auth.WithSecret("base64:czNjcjN0Cg=="),
		auth.WithHashAlgorithm("sha256"),
	)
	s.Require().Nil(err)

	s.e = e.(*auth.SecretEngine)
}

func (s *SecretAuthEngineTestSuite) TestSignAndVerify() {
	signature, err := s.e.Sign(map[string]interface{}{
		"a": "1",
		"b": 2,
		"c": []byte{3},
	})
	s.Require().Nil(err)

	sx := objx.New(signature)
	s.Equal("sha256", cast.ToString(sx.Get(auth.SECRET_CONTEXT_SIGNATURE_HASH_ALGORITHM).Inter()))
	s.Equal(int32(801072305), cast.ToInt32(sx.Get(auth.SECRET_CONTEXT_SIGNATURE_SESSION).Inter()))
	err = s.e.Verify(map[string]interface{}{
		"a": "1",
		"b": 2,
		"c": []byte{3},
	}, signature)
	s.Require().Nil(err)
}

func TestSecretAuthEngineTestSuite(t *testing.T) {
	suite.Run(t, new(SecretAuthEngineTestSuite))
}
