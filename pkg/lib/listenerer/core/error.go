package listenerer_core

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/lib/errors"
)

var (
	ErrListenerClosed         = fmt.Errorf("listener closed")
	ErrConnClosed             = fmt.Errorf("conn closed")
	ErrConnWaitEnabledTimeout = fmt.Errorf("conn wait enabled timeout")

	ErrUnsupportedNetwork, ErrUnsupportedNetworkFn = errors.NewErrorAndErrorFunc[string]("unsupported network")
)
