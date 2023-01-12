package transport_webrtc

import (
	"context"
)

func (t *WebrtcTransport) sendMessage(ctx context.Context, m Message) error {
	var n int

	logger := t.GetLogger().WithField("#method", "sendMessage").WithFields(t.wrapMessage(m))

	if t.rwc == nil {
		err := ErrInvalidSystemDataChannel
		logger.WithError(err).Debugf("system data channel is empty")
		return err
	}

	buf, err := t.marshaler.Marshal(&m)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal message")
		return err
	}

	if n, err = t.rwc.Write(buf); err != nil {
		logger.WithError(err).Debugf("failed to send buffer")
		return err
	}

	logger.WithField("bytes", n).Tracef("send message")

	return nil
}
