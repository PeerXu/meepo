package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) CloseTransport(target addr.Addr) error {
	return s.caller.Call(s.context(), sdk_core.METHOD_CLOSE_TRANSPORT, &sdk_interface.CloseTransportRequest{Target: target.String()}, nil)
}
