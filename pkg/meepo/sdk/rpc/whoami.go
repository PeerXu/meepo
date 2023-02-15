package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
)

func (s *RPCSDK) Whoami() (addr.Addr, error) {
	var a addr.Addr
	var addrStr string
	err := s.caller.Call(s.context(), sdk_core.METHOD_WHOAMI, rpc_core.NO_REQUEST(), &addrStr)
	if err != nil {
		return a, err
	}
	a, err = addr.FromString(addrStr)
	if err != nil {
		return a, err
	}
	return a, nil
}
