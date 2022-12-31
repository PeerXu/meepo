package sdk_rpc

import sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"

func (s *RPCSDK) ListChannels() (cvs []sdk_interface.ChannelView, err error) {
	err = s.caller.Call(s.context(), "listChannels", nil, &cvs)
	return
}
