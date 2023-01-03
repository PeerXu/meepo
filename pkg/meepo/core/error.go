package meepo_core

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/lib/errors"
)

var (
	ErrNoAvailableTrackers = fmt.Errorf("no available trackers")
)

var (
	ErrTeleportationNotFound, ErrTeleportationNotFoundFn = errors.NewErrorAndErrorFunc[string]("teleportation not found")
	ErrTransportNotFound, ErrTransportNotFoundFn         = errors.NewErrorAndErrorFunc[string]("transport not found")
	ErrTransportFound, ErrTransportFoundFn               = errors.NewErrorAndErrorFunc[string]("transport found")
	ErrTransportExist, ErrTransportExistFn               = errors.NewErrorAndErrorFunc[string]("transport exist")
	ErrTrackerNotFound, ErrTrackerNotFoundFn             = errors.NewErrorAndErrorFunc[string]("tracker not found")
	ErrInvalidNonce, ErrInvalidNonceFn                   = errors.NewErrorAndErrorFunc[uint32]("invalid nonce")
)
