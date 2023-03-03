package transport_webrtc

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/lib/errors"
)

var (
	ErrCallTimeout                = fmt.Errorf("call timeout")
	ErrInvalidPingSession         = fmt.Errorf("invalid ping session")
	ErrGatherTimeout              = fmt.Errorf("gather timeout")
	ErrInvalidAnswer              = fmt.Errorf("invalid answer")
	ErrRepeatedChannelID          = fmt.Errorf("repeated channel id")
	ErrInvalidSystemDataChannel   = fmt.Errorf("invalid system data channel")
	ErrNotConnectedPeerConnection = fmt.Errorf("not connected peer connection")
)

var (
	ErrSessionNotFound, ErrSessionNotFoundFn                 = errors.NewErrorAndErrorFunc[string]("session not found")
	ErrUnsupportedScope, ErrUnsupportedScopeFn               = errors.NewErrorAndErrorFunc[string]("unsupported scope")
	ErrUnsupportedMethod, ErrUnsupportedMethodFn             = errors.NewErrorAndErrorFunc[string]("unsupported method")
	ErrInvalidConnectionState, ErrInvalidConnectionStateFn   = errors.NewErrorAndErrorFunc[string]("invalid connection state")
	ErrChannelNotFound, ErrChannelNotFoundFn                 = errors.NewErrorAndErrorFunc[uint16]("channel not found")
	ErrPeerConnectionNotFound, ErrPeerConnectionNotFoundFn   = errors.NewErrorAndErrorFunc[Session]("peer connection not found")
	ErrPeerConnectionFound, ErrPeerConnectionFoundFn         = errors.NewErrorAndErrorFunc[Session]("peer connection found")
	ErrDataChannelNotFound, ErrDataChannelNotFoundFn         = errors.NewErrorAndErrorFunc[Session]("data channel not found")
	ErrReadWriteCloserNotFound, ErrReadWriteCloserNotFoundFn = errors.NewErrorAndErrorFunc[Session]("readWriteCloser not found")
)

func ErrInvalidPingSessionFn(expect int64, actual int64) error {
	return fmt.Errorf("%w: expect %v, actual %v", ErrInvalidPingSession, expect, actual)
}
