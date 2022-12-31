package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) GetChannel(target addr.Addr, id uint16) (cv sdk_interface.ChannelView, err error) {
	err = s.caller.Call(s.context(), "getChannel", &sdk_interface.GetChannelRequest{
		Target:    target.String(),
		ChannelID: id,
	}, &cv)
	return
}
