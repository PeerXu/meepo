package simple_sdk

import (
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	simple_logger "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/logger"
	"github.com/PeerXu/meepo/pkg/lib/config"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_simple_http "github.com/PeerXu/meepo/pkg/lib/rpc/simple_http"
	"github.com/PeerXu/meepo/pkg/meepo/sdk"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
)

func GetSDK() (sdk.SDK, error) {
	cfg := config.Get()

	logger, err := simple_logger.GetLogger()
	if err != nil {
		return nil, err
	}

	newSDKOpts := []sdk_core.NewSDKOption{
		well_known_option.WithLogger(logger),
	}
	newCallerOpts := []rpc_core.NewCallerOption{
		well_known_option.WithLogger(logger),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
	}
	switch cfg.Meepo.API.Name {
	case "http":
		name := "simple_http"
		newCallerOpts = append(newCallerOpts,
			rpc_simple_http.WithBaseURL("http://"+cfg.Meepo.API.Host),
		)
		caller, err := rpc.NewCaller(name, newCallerOpts...)
		if err != nil {
			return nil, err
		}
		newSDKOpts = append(newSDKOpts, rpc_core.WithCaller(caller))
	}

	return sdk.NewSDK("rpc", newSDKOpts...)
}
