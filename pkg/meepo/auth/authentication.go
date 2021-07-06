package auth

import "github.com/PeerXu/meepo/pkg/ofn"

type AuthenticateOption = ofn.OFN

type Authentication interface {
	Authenticate(sub string, opts ...AuthenticateOption) error
}
