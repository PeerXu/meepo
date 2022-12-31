package sdk_rpc

import sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"

func (s *RPCSDK) GetTeleportation(id string) (cv sdk_interface.TeleportationView, err error) {
	err = s.caller.Call(s.context(), "getTeleportation", &sdk_interface.GetTeleportationRequest{
		TeleportationID: id,
	}, &cv)
	return
}
