package transport_webrtc

import (
	"github.com/xtaci/smux"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) muxSessionAcceptLoop() {
	logger := t.GetLogger().WithField("#method", "muxSessionAcceptLoop")

	for {
		stm, err := t.muxSess.AcceptStream()
		if err != nil {
			logger.WithError(err).Debugf("failed to accept")
			return
		}
		go func(stm *smux.Stream) {
			label := t.parseMuxStreamLabel(stm)

			t.tempDataChannelsMtx.Lock()
			defer t.tempDataChannelsMtx.Unlock()

			tdc, found := t.tempDataChannels[label]
			if !found {
				t.tempDataChannels[label] = &tempDataChannel{upstream: stm}
				go t.removeTimeoutTempDataChannel(label)
				logger.Tracef("create temp data channel")
			} else {
				tdc.upstream = stm
				tdc.sinkChannel.setState(meepo_interface.CHANNEL_STATE_CONNECTING)
				go t.handleNewChannel(label, "muxSessionAcceptLoop")
			}
		}(stm)
	}
}
