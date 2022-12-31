package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) Ping(target addr.Addr, nonce uint32) (uint32, error) {
	var res sdk_interface.PingResponse
	err := s.caller.Call(s.context(), "ping", &sdk_interface.PingRequest{
		Target: target.String(),
		Nonce:  nonce,
	}, &res)
	if err != nil {
		return 0, err
	}

	return res.Nonce, nil
}
