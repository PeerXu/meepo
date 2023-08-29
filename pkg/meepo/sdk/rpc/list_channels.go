package sdk_rpc

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) ListChannels() (cvs []sdk_interface.ChannelView, err error) {
	err = s.caller.Call(s.context(), sdk_core.METHOD_LIST_CHANNELS, rpc_core.NO_REQUEST, &cvs)
	return
}
