package sdk_rpc

import (
	"github.com/PeerXu/meepo/pkg/internal/option"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
)

type RPCSDK struct {
	caller rpc_core.Caller
}

func NewRPCSDK(opts ...sdk_core.NewSDKOption) (sdk_core.SDK, error) {
	o := option.Apply(opts...)

	caller, err := rpc_core.GetCaller(o)
	if err != nil {
		return nil, err
	}

	return &RPCSDK{caller: caller}, nil
}

func init() {
	sdk_core.RegisterNewSDKFunc("rpc", NewRPCSDK)
}
