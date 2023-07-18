package sdk_rpc

import (
	"net"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) NewTeleportation(target addr.Addr, sourceAddr, sinkAddr net.Addr, mode string) (tpv sdk_interface.TeleportationView, err error) {
	err = s.caller.Call(s.context(), sdk_core.METHOD_NEW_TELEPORTATION, &sdk_interface.NewTeleportationRequest{
		Target: target.String(),
		TeleportationParam: sdk_interface.TeleportationParam{
			SourceNetwork: sourceAddr.Network(),
			SourceAddress: sourceAddr.String(),
			SinkNetwork:   sinkAddr.Network(),
			SinkAddress:   sinkAddr.String(),
			Mode:          mode,
		},
	}, &tpv)
	return
}