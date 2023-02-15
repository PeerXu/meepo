package sdk_rpc

import (
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) GetTeleportation(id string) (cv sdk_interface.TeleportationView, err error) {
	err = s.caller.Call(s.context(), sdk_core.METHOD_GET_TELEPORTATION, &sdk_interface.GetTeleportationRequest{
		TeleportationID: id,
	}, &cv)
	return
}
