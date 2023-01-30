package transport_webrtc

import (
	"context"
)

func (t *WebrtcTransport) sendMessage(ctx context.Context, m Message) error {
	var n int

	logger := t.GetLogger().WithField("#method", "sendMessage").WithFields(t.wrapMessage(m))

	rwc, err := t.loadSystemReadWriteCloserByContext(ctx)
	if err != nil {
		logger.WithError(err).Debugf("failed to get system rwc by context")
		return err
	}

	buf, err := t.marshaler.Marshal(&m)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal message")
		return err
	}

	if n, err = rwc.Write(buf); err != nil {
		logger.WithError(err).Debugf("failed to send buffer")
		return err
	}

	logger.WithField("bytes", n).Tracef("send message")

	return nil
}
