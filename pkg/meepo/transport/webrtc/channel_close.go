package transport_webrtc

import (
	"context"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (c *WebrtcSourceChannel) Close(context.Context) (err error) {
	logger := c.GetLogger().WithField("#method", "Close")

	if h := c.beforeCloseChannelHook; h != nil {
		if err = h(c, transport_core.WithIsSource(true)); err != nil {
			logger.WithError(err).Debugf("before close channel hook failed")
			return
		}
	}

	c.setState(meepo_interface.CHANNEL_STATE_CLOSING)

	if conn := c.Conn(); conn != nil {
		if err = conn.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close connection")
			return err
		}
	}

	if h := c.afterCloseChannelHook; h != nil {
		h(c, transport_core.WithIsSource(true))
	}

	c.setState(meepo_interface.CHANNEL_STATE_CLOSED)

	logger.Tracef("channel closed")

	return nil
}

func (c *WebrtcSinkChannel) Close(context.Context) (err error) {
	logger := c.GetLogger().WithField("#method", "Close")

	if h := c.beforeCloseChannelHook; h != nil {
		if err := h(c, transport_core.WithIsSink(true)); err != nil {
			logger.WithError(err).Debugf("before close channel hook failed")
			return err
		}
	}

	c.setState(meepo_interface.CHANNEL_STATE_CLOSING)

	if downstream := c.Conn(); downstream != nil {
		if err = downstream.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close conn")
			return err
		}
	}

	switch c.mode {
	case CHANNEL_MODE_RAW:
		if err = c.dc.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close webrtc.DataChannel")
			return err
		}
	case CHANNEL_MODE_MUX, CHANNEL_MODE_KCP:
		if err = c.upstream.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close ReadWriteCloser")
			return err
		}
	}

	if h := c.afterCloseChannelHook; h != nil {
		h(c, transport_core.WithIsSink(true))
	}

	c.setState(meepo_interface.CHANNEL_STATE_CLOSED)

	logger.Tracef("channel closed")

	return nil
}
