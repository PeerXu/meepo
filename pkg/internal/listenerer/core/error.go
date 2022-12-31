package listenerer_core

import "github.com/PeerXu/meepo/pkg/internal/errors"

var (
	ErrUnsupportedNetwork, ErrUnsupportedNetworkFn = errors.NewErrorAndErrorFunc[string]("unsupported network")
)
