package transport_webrtc

import (
	"context"
)

func (t *WebrtcTransport) doRequest(ctx context.Context, m Message) (err error) {
	logger := t.GetLogger().WithField("#method", "doRequest").WithFields(t.wrapMessage(m))

	t.polls.Store(t.parseResponseSession(m.Session), &LockableChannel{
		Ch: make(chan Message),
	})

	if err = t.sendMessage(ctx, m); err != nil {
		lch, ok := t.polls.LoadAndDelete(m.Session)
		if ok {
			lch.Close()
		}

		logger.WithError(err).Debugf("failed to send message")
		return err
	}

	logger.Tracef("do request")

	return nil
}
