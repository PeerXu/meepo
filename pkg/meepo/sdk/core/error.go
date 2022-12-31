package sdk_core

import "github.com/PeerXu/meepo/pkg/internal/errors"

var (
	ErrUnsupportedSDK, ErrUnsupportedSDKFn = errors.NewErrorAndErrorFunc[string]("unsupported sdk")
)
