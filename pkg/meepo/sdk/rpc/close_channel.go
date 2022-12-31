package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/lib/addr"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) CloseChannel(target addr.Addr, id uint16) error {
	return s.caller.Call(s.context(), "closeChannel", &sdk_interface.CloseChannelRequest{
		Target:    target.String(),
		ChannelID: id,
	}, rpc_core.NO_CONTENT())
}
