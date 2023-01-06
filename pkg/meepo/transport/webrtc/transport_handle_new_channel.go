package transport_webrtc

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
	mio "github.com/PeerXu/meepo/pkg/lib/io"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) handleNewChannel(label string) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "handleNewChannel",
		"label":   label,
	})

	t.tempDataChannelsMtx.Lock()
	tdc, ok := t.tempDataChannels[label]
	if !ok {
		t.tempDataChannelsMtx.Unlock()
		logger.Debugf("temp data channel not found")
		return
	}
	delete(t.tempDataChannels, label)
	t.tempDataChannelsMtx.Unlock()

	req := tdc.req
	logger = logger.WithFields(logging.Fields{
		"channelID": req.ChannelID,
		"network":   req.Network,
		"address":   req.Address,
		"mode":      req.Mode,
	})
	rwc := tdc.rwc

	ctx := t.context()

	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	var dc *webrtc.DataChannel
	switch req.Mode {
	case CHANNEL_MODE_RAW:
		dc = tdc.dc
	case CHANNEL_MODE_MUX:
		dc = t.muxDataChannel
	case CHANNEL_MODE_KCP:
		dc = t.kcpDataChannel
	default:
		logger.Debugf("unsupported mode")
		return
	}

	c := &WebrtcSinkChannel{
		WebrtcChannel: &WebrtcChannel{
			id:           req.ChannelID,
			sinkAddr:     dialer.NewAddr(req.Network, req.Address),
			logger:       t.GetRawLogger(),
			readyTimeout: t.readyTimeout,
			ready:        make(chan struct{}),
			mode:         req.Mode,
		},
		dc:  dc,
		rwc: rwc,
		onClose: func(c meepo_interface.Channel) error {
			t.csMtx.Lock()
			defer t.csMtx.Unlock()
			delete(t.cs, c.ID())
			return nil
		},
	}
	t.cs[c.ID()] = c

	go func() {
		var err error

		defer c.Close(ctx)

		dcConn := c.rwc
		conn, err := t.dialer.Dial(ctx, req.Network, req.Address)
		if err != nil {
			c.readyOnce.Do(func() {
				c.readyErr = err
				close(c.ready)
			})
			logger.WithError(err).Debugf("failed to dial")
			return
		}
		if c.dc.ReadyState() != webrtc.DataChannelStateOpen {
			defer conn.Close() // nolint:errcheck
			logger.Debugf("data channel closed")
			return
		}

		done1 := make(chan struct{})
		go func() {
			defer close(done1)
			n, err := mio.Copy(conn, dcConn)
			logger.WithError(err).WithFields(logging.Fields{
				"from":  "datachannel.ReadWriteCloser",
				"to":    "dialer.Conn",
				"bytes": n,
			}).Debugf("copy closed")
		}()

		done2 := make(chan struct{})
		go func() {
			defer close(done2)
			n, err := mio.Copy(dcConn, conn)
			logger.WithError(err).WithFields(logging.Fields{
				"from":  "dialer.Conn",
				"to":    "datachannel.ReadWriteCloser",
				"bytes": n,
			}).Debugf("copy closed")
		}()

		c.connVal.Store(conn)
		c.readyOnce.Do(func() { close(c.ready) })

		select {
		case <-done1:
		case <-done2:
		}
	}()

	logger.Tracef("handle new channel")
}
