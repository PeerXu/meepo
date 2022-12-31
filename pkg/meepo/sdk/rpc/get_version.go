package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/version"
)

func (s *RPCSDK) GetVersion() (*version.V, error) {
	var v version.V
	err := s.caller.Call(s.context(), "getVersion", nil, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
