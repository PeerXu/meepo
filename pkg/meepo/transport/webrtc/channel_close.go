package transport_webrtc

import "context"

func (c *WebrtcSourceChannel) Close(context.Context) (err error) {
	logger := c.GetLogger().WithField("#method", "Close")

	if c.onClose != nil {
		if err = c.onClose(c); err != nil {
			logger.WithError(err).Debugf("failed to onClose")
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
		if err = c.rwc.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close ReadWriteCloser")
			return err
		}
	}

	logger.Tracef("channel closed")

	return nil
}

func (c *WebrtcSinkChannel) Close(context.Context) (err error) {
	logger := c.GetLogger().WithField("#method", "Close")

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

	switch c.mode {
	case CHANNEL_MODE_RAW:
		if err = c.dc.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close webrtc.DataChannel")
			return err
		}
	case CHANNEL_MODE_MUX, CHANNEL_MODE_KCP:
		if err = c.rwc.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close ReadWriteCloser")
			return err
		}
	}

	logger.Tracef("channel closed")

	return nil
}
