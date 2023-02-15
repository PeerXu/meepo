package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) CloseChannel(target addr.Addr, id uint16) error {
	return s.caller.Call(s.context(), sdk_core.METHOD_CLOSE_CHANNEL, &sdk_interface.CloseChannelRequest{
		Target:    target.String(),
		ChannelID: id,
	}, nil)
}
