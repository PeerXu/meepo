package meepo_socks5

import "github.com/PeerXu/meepo/pkg/internal/errors"

var (
	ErrInvalidDomain, ErrInvalidDomainFn = errors.NewErrorAndErrorFunc[string]("invalid domain")
)
