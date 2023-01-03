package transport_pipe

import (
	"context"

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

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "NewChannel",
		"network": network,
		"address": address,
	})

	readyTimeout, err := transport_core.GetReadyTimeout(o)
	if err != nil {
		return nil, err
	}

	channelID := t.nextChannelID()
	logger = logger.WithField("channelID", channelID)

	c := &PipeChannel{
		id:       channelID,
		sinkAddr: dialer.NewAddr(network, address),
		onClose: func(c meepo_interface.Channel) error {
			t.csMtx.Lock()
			defer t.csMtx.Unlock()
			delete(t.cs, c.ID())
			return nil
		},
		logger:       t.GetRawLogger().WithField("addr", t.Addr()),
		readyTimeout: readyTimeout,
		ready:        make(chan struct{}),
	}
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

	logger.Tracef("new channel")

	return c, nil
}
