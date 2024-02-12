package transport_webrtc

import (
	"context"

	matomic "github.com/PeerXu/meepo/pkg/lib/atomic"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type NewChannelRequest struct {
	Label     string
	ChannelID uint16
	Network   string
	Address   string
	Mode      string
}

func (t *WebrtcTransport) newRemoteChannel(ctx context.Context, mode string, label string, channelID uint16, network, address string) error {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":   "newRemoteChannel",
		"label":     label,
		"channelID": channelID,
		"network":   network,
		"address":   address,
		"mode":      mode,
	})

	if err := t.Call(ctx, SYS_METHOD_NEW_CHANNEL, &NewChannelRequest{
		Label:     label,
		ChannelID: channelID,
		Network:   network,
		Address:   address,
		Mode:      mode,
	}, nil, well_known_option.WithScope("sys")); err != nil {
		logger.WithError(err).Debugf("failed to new remote channel")
		return err
	}

	logger.Tracef("new remote channel")

	return nil
}

func (t *WebrtcTransport) onNewChannel(ctx context.Context, req NewChannelRequest) (res any, err error) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":   "onNewChannel",
		"label":     req.Label,
		"channelID": req.ChannelID,
		"network":   req.Network,
		"address":   req.Address,
		"mode":      req.Mode,
	})

	t.tempDataChannelsMtx.Lock()
	defer t.tempDataChannelsMtx.Unlock()

	if h := t.BeforeNewChannelHook; h != nil {
		if err = h(req.Network, req.Address, transport_core.WithIsSink(true)); err != nil {
			logger.WithError(err).Debugf("before new channel hook failed")
			return nil, err
		}
	}

	if _, found := t.cs[req.ChannelID]; found {
		err = ErrRepeatedChannelID
		logger.WithError(err).Debugf("repeated channel id")
		return
	}

	sinkChannel := &WebrtcSinkChannel{
		// upstream required
		WebrtcChannel: &WebrtcChannel{
			// dc required
			id:            req.ChannelID,
			sinkAddr:      dialer.NewAddr(req.Network, req.Address),
			logger:        t.GetRawLogger(),
			s:             matomic.NewValue[meepo_interface.ChannelState](),
			readyTimeout:  t.readyTimeout,
			readyCh:       make(chan struct{}),
			mode:          req.Mode,
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
		downstreamVal: matomic.NewValue[meepo_interface.Conn](),
	}

	t.addChannel(sinkChannel)
	sinkChannel.setState(meepo_interface.CHANNEL_STATE_NEW)

	tdc, found := t.tempDataChannels[req.Label]
	if !found {

		t.tempDataChannels[req.Label] = &tempDataChannel{
			request:     &req,
			sinkChannel: sinkChannel,
		}

		go t.removeTimeoutTempDataChannel(req.Label)
		logger.Tracef("create temp data channel")
		return
	}

	tdc.request = &req
	tdc.sinkChannel = sinkChannel

	if tdc.upstream == nil {
		logger.Tracef("wait for data channel open")
		return
	}

	sinkChannel.setState(meepo_interface.CHANNEL_STATE_CONNECTING)
	go t.handleNewChannel(req.Label, "onNewChannel")

	logger.Tracef("on new channel")

	return
}
