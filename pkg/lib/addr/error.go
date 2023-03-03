package addr

import (
	"github.com/PeerXu/meepo/pkg/lib/errors"
)

var (
	ErrInvalidAddrString, ErrInvalidAddrStringFn = errors.NewErrorAndErrorFunc[string]("invalid addr string")
)
