package listenerer_core

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/lib/errors"
)

var (
	ErrListenerClosed = fmt.Errorf("listener closed")

	ErrUnsupportedNetwork, ErrUnsupportedNetworkFn = errors.NewErrorAndErrorFunc[string]("unsupported network")
)
