package sdk_rpc

import (
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) CloseTeleportation(id string) error {
	return s.caller.Call(s.context(), sdk_core.METHOD_CLOSE_TELEPORTATION, &sdk_interface.CloseTeleportationRequest{
		TeleportationID: id,
	}, nil)
}
