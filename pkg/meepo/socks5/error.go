package meepo_socks5

import "github.com/PeerXu/meepo/pkg/lib/errors"

var (
	ErrInvalidDomain, ErrInvalidDomainFn = errors.NewErrorAndErrorFunc[string]("invalid domain")
)
