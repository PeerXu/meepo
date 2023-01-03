//go:build js

package meepo_core

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

func newWebrtcSettingEngine(o option.Option) (se webrtc.SettingEngine, err error) {
	se.DetachDataChannels()
	return
}
