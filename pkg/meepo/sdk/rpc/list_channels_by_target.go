package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) ListChannelsByTarget(target addr.Addr) (cvs []sdk_interface.ChannelView, err error) {
	err = s.caller.Call(s.context(), sdk_core.METHOD_LIST_CHANNELS_BY_TARGET, &sdk_interface.ListChannelsByTarget{
		Target: target.String(),
	}, &cvs)
	return
}
