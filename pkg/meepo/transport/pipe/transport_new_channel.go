package transport_pipe

import (
	"context"

	matomic "github.com/PeerXu/meepo/pkg/lib/atomic"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *PipeTransport) NewChannel(ctx context.Context, network string, address string, opts ...meepo_interface.NewChannelOption) (meepo_interface.Channel, error) {
	o := option.ApplyWithDefault(defaultNewPipeTransportOptions(), opts...)

	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	channelID := t.nextChannelID()
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":   "NewChannel",
		"network":   network,
		"address":   address,
		"channelID": channelID,
	})

	if h := t.BeforeNewChannelHook; h != nil {
		if err := h(
			network, address,
			transport_core.WithIsSource(true),
			transport_core.WithIsSink(true),
		); err != nil {
			logger.WithError(err).Debugf("before new channel hook failed")
			return nil, err
		}
	}

	readyTimeout, err := transport_core.GetReadyTimeout(o)
	if err != nil {
		return nil, err
	}

	c := &PipeChannel{
		id:            channelID,
		state:         matomic.NewValue[meepo_interface.ChannelState](),
		sinkAddr:      dialer.NewAddr(network, address),
		logger:        t.GetRawLogger().WithField("addr", t.Addr()),
		onStateChange: t.onChannelStateChange,
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
		readyTimeout:          readyTimeout,
		ready:                 make(chan struct{}),
	}
	c.setState(meepo_interface.CHANNEL_STATE_NEW)
	t.cs[c.ID()] = c
	go func() {
		defer c.readyOnce.Do(func() { close(c.ready) })

		c.setState(meepo_interface.CHANNEL_STATE_CONNECTING)

		conn, err := t.dialer.Dial(ctx, network, address)
		if err != nil {
			c.readyErr = err
			c.setState(meepo_interface.CHANNEL_STATE_CLOSING)
			go c.Close(c.context())

			logger.WithError(err).Debugf("failed to dial")
			return
		}

		c.conn = conn
		c.setState(meepo_interface.CHANNEL_STATE_OPEN)
	}()

	if h := t.AfterNewChannelHook; h != nil {
		h(c, transport_core.WithIsSource(true), transport_core.WithIsSink(true))
	}

	logger.Tracef("new channel")

	return c, nil
}
