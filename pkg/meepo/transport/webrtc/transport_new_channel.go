package transport_webrtc

import (
	"context"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
	dialer_interface "github.com/PeerXu/meepo/pkg/lib/dialer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) NewChannel(ctx context.Context, network string, address string, opts ...meepo_interface.NewChannelOption) (meepo_interface.Channel, error) {
	o := option.ApplyWithDefault(defaultNewChannelOptions(), opts...)

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "NewChannel",
		"network": network,
		"address": address,
	})

	mode, err := well_known_option.GetMode(o)
	if err != nil {
		logger.WithError(err).Debugf("failed to get mode")
		return nil, err
	}
	logger = logger.WithField("mode", mode)

	var label string
	var dc *webrtc.DataChannel
	var rwc dialer_interface.Conn
	var closer func() error
	channelID := t.nextChannelID()
	logger = logger.WithField("channelID", channelID)

	switch mode {
	case CHANNEL_MODE_RAW:
		label = t.nextDataChannelLabel("data")
		dc, err = t.pc.CreateDataChannel(label, nil)
		if err != nil {
			logger.WithError(err).Debugf("failed to create data channel")
			return nil, err
		}
		closer = dc.Close
	case CHANNEL_MODE_MUX:
		stm, err := t.muxSess.OpenStream()
		if err != nil {
			logger.WithError(err).Debugf("failed to create mux stream")
			return nil, err
		}
		label = t.parseMuxStreamLabel(stm)
		dc = t.muxDataChannel
		rwc = stm
		closer = stm.Close
	case CHANNEL_MODE_KCP:
		stm, err := t.kcpSess.OpenStream()
		if err != nil {
			logger.WithError(err).Debugf("failed to create kcp stream")
			return nil, err
		}
		label = t.parseKcpStreamLabel(stm)
		dc = t.kcpDataChannel
		rwc = stm
		closer = stm.Close
	}
	logger = logger.WithField("label", label)

	err = t.newRemoteChannel(ctx, mode, label, channelID, network, address)
	if err != nil {
		defer closer() // nolint:errcheck
		logger.WithError(err).Debugf("failed to new remote channel")
		return nil, err
	}

	c := &WebrtcSourceChannel{
		WebrtcChannel: &WebrtcChannel{
			id:           channelID,
			sinkAddr:     dialer.NewAddr(network, address),
			logger:       t.GetRawLogger(),
			readyTimeout: t.readyTimeout,
			ready:        make(chan struct{}),
			mode:         mode,
		},
		dc: dc,
		onClose: func(c meepo_interface.Channel) error {
			t.csMtx.Lock()
			defer t.csMtx.Unlock()
			delete(t.cs, c.ID())
			return nil
		},
	}

	// TODO: channel done
	switch mode {
	case CHANNEL_MODE_RAW:
		t.afterRawSourceChannelCreate(dc, c)
	case CHANNEL_MODE_MUX:
		t.csMtx.Lock()
		defer t.csMtx.Unlock()
		c.rwc = rwc
		c.conn = rwc
		t.cs[channelID] = c
		c.readyOnce.Do(func() { close(c.ready) })
	case CHANNEL_MODE_KCP:
		t.csMtx.Lock()
		defer t.csMtx.Unlock()
		c.rwc = rwc
		c.conn = rwc
		t.cs[channelID] = c
		c.readyOnce.Do(func() { close(c.ready) })
	}

	logger.Tracef("new channel")

	return c, nil
}
