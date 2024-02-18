package meepo_debug_sdk

import (
	meepo_debug_sdk_core "github.com/PeerXu/meepo/pkg/meepo/debug/sdk/core"
	_ "github.com/PeerXu/meepo/pkg/meepo/debug/sdk/http"
)

var NewSDK = meepo_debug_sdk_core.New
