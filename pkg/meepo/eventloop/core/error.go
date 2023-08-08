package meepo_eventloop_core

import "github.com/PeerXu/meepo/pkg/lib/errors"

var (
	ErrHandlerExisted, ErrHandlerExistedFn = errors.NewErrorAndErrorFunc[string]("handler existed")
)
