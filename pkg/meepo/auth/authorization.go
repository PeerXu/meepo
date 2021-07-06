package auth

import "github.com/PeerXu/meepo/pkg/ofn"

type AuthorizeOption = ofn.OFN

type Authorization interface {
	Authorize(sub, obj, act string, opts ...AuthorizeOption) error
}
