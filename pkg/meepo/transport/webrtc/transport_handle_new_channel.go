package transport_webrtc

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
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
		beforeCloseChannelHook: func(c meepo_interface.Channel, opts ...transport_core.HookOption) error {
			if h := t.BeforeCloseChannelHook; h != nil {
				if err := h(c, opts...); err != nil {
					return err
				}
			}

			t.csMtx.Lock()
			defer t.csMtx.Unlock()
			delete(t.cs, c.ID())

			return nil

		},
		afterCloseChannelHook: t.AfterCloseChannelHook,
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

	if h := t.AfterNewChannelHook; h != nil {
		h(c, transport_core.WithIsSink(true))
	}

	logger.Tracef("handle new channel")
}
