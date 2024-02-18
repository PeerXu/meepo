package lib_registerer

import "github.com/PeerXu/meepo/pkg/lib/errors"

var (
	ErrUnsupported, ErrUnsupportedFn = errors.NewErrorAndErrorFunc[string]("unsupported")
)
