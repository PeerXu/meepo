package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) CloseTransport(target addr.Addr) error {
	return s.caller.Call(s.context(), "closeTransport", &sdk_interface.CloseTransportRequest{Target: target.String()}, rpc_core.NO_CONTENT())
}
