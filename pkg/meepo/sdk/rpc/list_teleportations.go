package sdk_rpc

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) ListTeleportations() (tpvs []sdk_interface.TeleportationView, err error) {
	err = s.caller.Call(s.context(), "listTeleportations", rpc_core.NO_CONTENT(), &tpvs)
	return
}
