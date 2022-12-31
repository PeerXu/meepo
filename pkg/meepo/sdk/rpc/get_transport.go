package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) GetTransport(target addr.Addr) (tv sdk_interface.TransportView, err error) {
	err = s.caller.Call(s.context(), "getTransport", &sdk_interface.GetTransportRequest{
		Target: target.String(),
	}, &tv)
	return
}
