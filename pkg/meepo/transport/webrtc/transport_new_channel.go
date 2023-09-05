package transport_webrtc

import (
	"context"

	matomic "github.com/PeerXu/meepo/pkg/lib/atomic"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) NewChannel(ctx context.Context, network string, address string, opts ...meepo_interface.NewChannelOption) (meepo_interface.Channel, error) {
	o := option.ApplyWithDefault(defaultNewChannelOptions(), opts...)

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "NewChannel",
		"network": network,
		"address": address,
	})

	if h := t.BeforeNewChannelHook; h != nil {
		if err := h(network, address, transport_core.WithIsSource(true)); err != nil {
			logger.WithError(err).Debugf("before new channel hook failed")
			return nil, err
		}
	}

	mode, err := well_known_option.GetMode(o)
	if err != nil {
		logger.WithError(err).Debugf("failed to get mode")
		return nil, err
	}
	logger = logger.WithField("mode", mode)

	var label string
	channelID := t.nextChannelID()
	logger = logger.WithField("channelID", channelID)

	c := &WebrtcSourceChannel{
		WebrtcChannel: &WebrtcChannel{
			id:            channelID,
			sinkAddr:      dialer.NewAddr(network, address),
			logger:        t.GetRawLogger(),
			s:             matomic.NewValue[meepo_interface.ChannelState](),
			readyTimeout:  t.readyTimeout,
			readyCh:       make(chan struct{}),
			mode:          mode,
			onStateChange: t.onChannelStateChange,
			beforeCloseChannelHook: func(c meepo_interface.Channel, opts ...transport_core.HookOption) error {
				if h := t.BeforeCloseChannelHook; h != nil {
					if err := h(c, opts...); err != nil {
						return err
					}
				}

				t.removeChannel(c)

				return nil
			},
			afterCloseChannelHook: t.AfterCloseChannelHook,
		},
	}
	t.addChannel(c)
	c.setState(meepo_interface.CHANNEL_STATE_NEW)

	switch mode {
	case CHANNEL_MODE_RAW:
		pc, err := t.loadRandomPeerConnection()
		if err != nil {
			logger.WithError(err).Debugf("failed to load peer connection")
			return nil, err
		}

		// assign c.conn when webrtc datachannel open
		label = t.nextDataChannelLabel("data")
		c.dc, err = pc.CreateDataChannel(label, nil)
		if err != nil {
			defer c.Close(ctx)
			logger.WithError(err).Debugf("failed to create data channel")
			return nil, err
		}
	case CHANNEL_MODE_MUX:
		stm, err := t.muxSess.OpenStream()
		if err != nil {
			defer c.Close(ctx)
			logger.WithError(err).Debugf("failed to create mux stream")
			return nil, err
		}
		label = t.parseMuxStreamLabel(stm)
		c.dc = t.muxDataChannel
		c.conn = stm
	case CHANNEL_MODE_KCP:
		stm, err := t.kcpSess.OpenStream()
		if err != nil {
			defer c.Close(ctx)
			logger.WithError(err).Debugf("failed to create kcp stream")
			return nil, err
		}
		label = t.parseKcpStreamLabel(stm)
		c.dc = t.kcpDataChannel
		c.conn = stm
	}
	logger = logger.WithField("label", label)

	c.setState(meepo_interface.CHANNEL_STATE_CONNECTING)
	err = t.newRemoteChannel(ctx, mode, label, channelID, network, address)
	if err != nil {
		defer c.Close(ctx) // nolint:errcheck
		logger.WithError(err).Debugf("failed to new remote channel")
		return nil, err
	}

	switch mode {
	case CHANNEL_MODE_RAW:
		t.afterRawSourceChannelCreate(c.dc, c)
	case CHANNEL_MODE_MUX, CHANNEL_MODE_KCP:
		c.setState(meepo_interface.CHANNEL_STATE_OPEN)
		c.ready()
	}

	if h := t.AfterNewChannelHook; h != nil {
		h(c, transport_core.WithIsSource(true))
	}

	logger.Tracef("new channel")

	return c, nil
}
