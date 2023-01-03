package transport_core

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/lib/errors"
)

var (
	ErrReadyTimeout    = fmt.Errorf("ready timeout")
	ErrTransportClosed = fmt.Errorf("transport closed")

	ErrUnsupportedTransport, ErrUnsupportedTransportFn = errors.NewErrorAndErrorFunc[string]("unsupported transport")
	ErrChannelNotFound, ErrChannelNotFoundFn           = errors.NewErrorAndErrorFunc[uint16]("channel not found")
	ErrUnsupportedMethod, ErrUnsupportedMethodFn       = errors.NewErrorAndErrorFunc[string]("unsupported method")
)
