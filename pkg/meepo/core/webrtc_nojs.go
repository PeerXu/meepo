//go:build !js

package meepo_core

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

func newWebrtcSettingEngine(o option.Option) (se webrtc.SettingEngine, err error) {
	recvBufSize, err := well_known_option.GetWebrtcReceiveBufferSize(o)
	if err != nil {
		return
	}

	se.DetachDataChannels()
	se.SetSCTPMaxReceiveBufferSize(recvBufSize)

	return
}
