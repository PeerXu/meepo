package sdk_rpc

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) ListTransports() (ts []sdk_interface.TransportView, err error) {
	err = s.caller.Call(s.context(), sdk_core.METHOD_LIST_TRANSPORTS, rpc_core.NO_REQUEST(), &ts)
	return
}
