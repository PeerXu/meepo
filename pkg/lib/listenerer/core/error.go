package listenerer_core

import "github.com/PeerXu/meepo/pkg/lib/errors"

var (
	ErrUnsupportedNetwork, ErrUnsupportedNetworkFn = errors.NewErrorAndErrorFunc[string]("unsupported network")
)
