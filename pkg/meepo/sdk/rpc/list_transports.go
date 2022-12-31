package sdk_rpc

import sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"

func (s *RPCSDK) ListTransports() (ts []sdk_interface.TransportView, err error) {
	err = s.caller.Call(s.context(), "listTransports", nil, &ts)
	return
}
