package transport_webrtc

import (
	"github.com/pion/webrtc/v3"

	mio "github.com/PeerXu/meepo/pkg/lib/io"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) handleNewChannel(label string, fromMethod string) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":     "handleNewChannel",
		"label":       label,
		"#fromMethod": fromMethod,
	})

	ctx := t.context()

	t.tempDataChannelsMtx.Lock()
	tdc, ok := t.tempDataChannels[label]
	if !ok {
		t.tempDataChannelsMtx.Unlock()
		logger.Debugf("temp data channel not found")
		return
	}
	delete(t.tempDataChannels, label)
	t.tempDataChannelsMtx.Unlock()

	req := tdc.request
	logger = logger.WithFields(logging.Fields{
		"channelID": req.ChannelID,
		"network":   req.Network,
		"address":   req.Address,
		"mode":      req.Mode,
	})
	sinkChannel := tdc.sinkChannel
	sinkChannel.upstream = tdc.upstream
	switch req.Mode {
	case CHANNEL_MODE_RAW:
		sinkChannel.dc = tdc.datachannel
	case CHANNEL_MODE_MUX:
		sinkChannel.dc = t.muxDataChannel
	case CHANNEL_MODE_KCP:
		sinkChannel.dc = t.kcpDataChannel
	default:
		if err := sinkChannel.Close(ctx); err != nil {
			logger.WithError(err).Debugf("failed to close sink channel")
		}
		logger.Debugf("unsupported mode")
		return
	}

	go func(c *WebrtcSinkChannel) {
		defer c.Close(ctx)

		upstream := c.upstream
		downstream, err := t.dialer.Dial(ctx, req.Network, req.Address)
		if err != nil {
			c.readyWithError(err)
			logger.WithError(err).Debugf("failed to dial")
			return
		}
		defer downstream.Close() // nolint:errcheck

		if c.dc.ReadyState() != webrtc.DataChannelStateOpen {
			logger.Debugf("data channel closed")
			return
		}

		done1 := make(chan struct{})
		go func() {
			defer close(done1)
			n, err := mio.Copy(downstream, upstream)
			logger.WithError(err).WithFields(logging.Fields{
				"from":  "datachannel.ReadWriteCloser",
				"to":    "dialer.Conn",
				"bytes": n,
			}).Debugf("copy closed")
		}()

		done2 := make(chan struct{})
		go func() {
			defer close(done2)
			n, err := mio.Copy(upstream, downstream)
			logger.WithError(err).WithFields(logging.Fields{
				"from":  "dialer.Conn",
				"to":    "datachannel.ReadWriteCloser",
				"bytes": n,
			}).Debugf("copy closed")
		}()

		c.downstreamVal.Store(downstream)
		c.setState(meepo_interface.CHANNEL_STATE_OPEN)
		c.ready()

		select {
		case <-done1:
		case <-done2:
		}
	}(sinkChannel)

	if h := t.AfterNewChannelHook; h != nil {
		h(sinkChannel, transport_core.WithIsSink(true))
	}

	logger.Tracef("handle new channel")
}
