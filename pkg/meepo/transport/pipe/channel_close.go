package transport_pipe

import (
	"context"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (c *PipeChannel) Close(context.Context) (err error) {
	logger := c.GetLogger().WithField("#method", "Close")

	if h := c.beforeCloseChannelHook; h != nil {
		if err := h(c, transport_core.WithIsSource(true), transport_core.WithIsSink(true)); err != nil {
			logger.WithError(err).Debugf("before close channel hook failed")
			return err
		}
	}

	c.setState(meepo_interface.CHANNEL_STATE_CLOSING)

	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close conn")
			return err
		}
	}

	if h := c.afterCloseChannelHook; h != nil {
		h(c, transport_core.WithIsSource(true), transport_core.WithIsSink(true))
	}

	c.setState(meepo_interface.CHANNEL_STATE_CLOSED)

	logger.Tracef("channel closed")

	return nil
}
