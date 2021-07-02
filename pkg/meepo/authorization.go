package meepo

import (
	"encoding/base64"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"
	"golang.org/x/crypto/bcrypt"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
)

func (mp *Meepo) GetAuthorizationName() string {
	return cast.ToString(mp.opt.Get("authorizationName").Inter())
}

func (mp *Meepo) Authorize(subject, object, action string, opts ...auth.AuthorizeOption) error {
	authName := mp.GetAuthorizationName()
	switch authName {
	case "secret":
		return mp.secretAuthorize(subject, object, action, opts...)
	case "dummy":
		fallthrough
	default:
		return nil
	}
}

func (mp *Meepo) secretAuthorize(subject, object, action string, opts ...auth.AuthorizeOption) error {
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "secretAuthorize",
		"subject": subject,
		"object":  object,
		"action":  action,
	})

	o := objx.New(map[string]interface{}{})
	for _, opt := range opts {
		opt(o)
	}

	switch action {
	case string(METHOD_NEW_TELEPORTATION):
		hashedSecret, err := base64.StdEncoding.DecodeString(cast.ToString(o.Get("authorizationSecret").Inter()))
		if err != nil {
			return fmt.Errorf("%w: invalid hashed secret", ErrUnauthorized)
		}

		if len(hashedSecret) == 0 {
			return fmt.Errorf("%w: require secret", ErrUnauthorized)
		}

		secret := cast.ToString(mp.opt.Get("authorizationSecret").Inter())

		if err = bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret)); err != nil {
			logger.WithError(err).Debugf("incorrect password")
			return fmt.Errorf("%w: incorrect password", ErrUnauthorized)
		}
	default:
	}

	return nil
}
