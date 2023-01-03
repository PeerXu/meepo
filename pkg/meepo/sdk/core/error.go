package sdk_core

import "github.com/PeerXu/meepo/pkg/lib/errors"

var (
	ErrUnsupportedSDK, ErrUnsupportedSDKFn = errors.NewErrorAndErrorFunc[string]("unsupported sdk")
)
