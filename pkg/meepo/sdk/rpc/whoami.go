package sdk_rpc

import "github.com/PeerXu/meepo/pkg/lib/addr"

func (s *RPCSDK) Whoami() (addr.Addr, error) {
	var a addr.Addr
	var addrStr string
	err := s.caller.Call(s.context(), "whoami", nil, &addrStr)
	if err != nil {
		return a, err
	}
	a, err = addr.FromString(addrStr)
	if err != nil {
		return a, err
	}
	return a, nil
}
