package sdk

import (
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	_ "github.com/PeerXu/meepo/pkg/meepo/sdk/rpc"
)

type SDK = sdk_core.SDK

var NewSDK = sdk_core.NewSDK
