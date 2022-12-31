package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) NewTransport(target addr.Addr, opts ...sdk_interface.NewTransportOption) (sdk_interface.TransportView, error) {
	o := option.Apply(opts...)

	var req sdk_interface.NewTransportRequest
	req.Target = target.String()
	req.Manual, _ = well_known_option.GetManual(o)
	if req.Manual {
		sp := &req.SmuxParam
		sp.EnableMux, _ = well_known_option.GetEnableMux(o)
		if sp.EnableMux {
			sp.MuxVer, _ = well_known_option.GetMuxVer(o)
			sp.MuxBuf, _ = well_known_option.GetMuxBuf(o)
			sp.MuxStreamBuf, _ = well_known_option.GetMuxStreamBuf(o)
			sp.MuxNocomp, _ = well_known_option.GetMuxNocomp(o)
		}

		kp := &req.KcpParam
		kp.EnableKcp, _ = well_known_option.GetEnableKcp(o)
		if kp.EnableKcp {
			kp.KcpPreset, _ = well_known_option.GetKcpPreset(o)
			kp.KcpCrypt, _ = well_known_option.GetKcpCrypt(o)
			kp.KcpKey, _ = well_known_option.GetKcpKey(o)
			kp.KcpMtu, _ = well_known_option.GetKcpMtu(o)
			kp.KcpSndwnd, _ = well_known_option.GetKcpSndwnd(o)
			kp.KcpRcvwnd, _ = well_known_option.GetKcpRcvwnd(o)
			kp.KcpDataShard, _ = well_known_option.GetKcpDataShard(o)
			kp.KcpParityShard, _ = well_known_option.GetKcpParityShard(o)
		}
	}

	var tv sdk_interface.TransportView
	err := s.caller.Call(s.context(), "newTransport", &req, &tv)
	if err != nil {
		return tv, err
	}
	return tv, nil
}
