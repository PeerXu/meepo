package sdk_rpc

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	"github.com/PeerXu/meepo/pkg/lib/version"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
)

func (s *RPCSDK) GetVersion() (*version.V, error) {
	var v version.V
	err := s.caller.Call(s.context(), sdk_core.METHOD_GET_VERSION, rpc_core.NO_REQUEST, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
