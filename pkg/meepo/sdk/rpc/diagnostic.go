package sdk_rpc

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) Diagnostic() (sdk_interface.DiagnosticReportView, error) {
	res := make(sdk_interface.DiagnosticReportView)
	err := s.caller.Call(s.context(), sdk_core.METHOD_DIAGNOSTIC, rpc_core.NO_REQUEST, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
