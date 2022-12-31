package sdk_rpc

import (
	"net"

	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) Teleport(target addr.Addr, sourceAddr, sinkAddr net.Addr, mode string, opts ...sdk_interface.TeleportOption) (sdk_interface.TeleportationView, error) {
	o := option.Apply(opts...)

	var req sdk_interface.TeleportRequest
	req.Target = target.String()

	req.Manual, _ = well_known_option.GetManual(o)
	if req.Manual {
		req.EnableMux, _ = well_known_option.GetEnableMux(o)
		if req.EnableMux {
			req.MuxVer, _ = well_known_option.GetMuxVer(o)
			req.MuxBuf, _ = well_known_option.GetMuxBuf(o)
			req.MuxStreamBuf, _ = well_known_option.GetMuxStreamBuf(o)
			req.MuxNocomp, _ = well_known_option.GetMuxNocomp(o)
		}

		req.EnableKcp, _ = well_known_option.GetEnableKcp(o)
		if req.EnableKcp {
			req.KcpPreset, _ = well_known_option.GetKcpPreset(o)
			req.KcpCrypt, _ = well_known_option.GetKcpCrypt(o)
			req.KcpKey, _ = well_known_option.GetKcpKey(o)
			req.KcpMtu, _ = well_known_option.GetKcpMtu(o)
			req.KcpSndwnd, _ = well_known_option.GetKcpSndwnd(o)
			req.KcpRcvwnd, _ = well_known_option.GetKcpRcvwnd(o)
			req.KcpDataShard, _ = well_known_option.GetKcpDataShard(o)
			req.KcpParityShard, _ = well_known_option.GetKcpParityShard(o)
		}
	}

	req.Mode = mode
	req.SourceNetwork = sourceAddr.Network()
	req.SourceAddress = sourceAddr.String()
	req.SinkNetwork = sinkAddr.Network()
	req.SinkAddress = sinkAddr.String()

	var tpv sdk_interface.TeleportationView
	err := s.caller.Call(s.context(), "teleport", &req, &tpv)
	if err != nil {
		return tpv, err
	}

	return tpv, nil
}
