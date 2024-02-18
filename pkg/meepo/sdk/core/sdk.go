package sdk_core

import (
	lib_registerer "github.com/PeerXu/meepo/pkg/lib/registerer"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

type SDK = sdk_interface.SDK

var RegisterNewSDKFunc, NewSDK = lib_registerer.Pair[SDK]()
