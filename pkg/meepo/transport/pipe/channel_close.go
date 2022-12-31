package transport_pipe

import (
	"context"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (c *PipeChannel) Close(context.Context) (err error) {
	logger := c.GetLogger().WithField("#method", "Close")

	defer c.setState(meepo_interface.CHANNEL_STATE_CLOSED)

	c.setState(meepo_interface.CHANNEL_STATE_CLOSING)
	if c.onClose != nil {
		if err = c.onClose(c); err != nil {
			logger.WithError(err).Debugf("failed to onClose")
			return err
		}
	}

	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close conn")
			return err
		}
	}

	logger.Tracef("channel closed")

	return nil
}
