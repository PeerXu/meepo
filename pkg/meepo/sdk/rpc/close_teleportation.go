package sdk_rpc

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) CloseTeleportation(id string) error {
	return s.caller.Call(s.context(), "closeTeleportation", &sdk_interface.CloseTeleportationRequest{
		TeleportationID: id,
	}, rpc_core.NO_CONTENT())
}
