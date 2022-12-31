package sdk_core

import (
	"sync"

	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

type SDK = sdk_interface.SDK

type NewSDKFunc func(...NewSDKOption) (SDK, error)

var newSDKFuncs sync.Map

func NewSDK(name string, opts ...NewSDKOption) (SDK, error) {
	v, ok := newSDKFuncs.Load(name)
	if !ok {
		return nil, ErrUnsupportedSDKFn(name)
	}
	return v.(NewSDKFunc)(opts...)
}

func RegisterNewSDKFunc(name string, fn NewSDKFunc) {
	newSDKFuncs.Store(name, fn)
}
