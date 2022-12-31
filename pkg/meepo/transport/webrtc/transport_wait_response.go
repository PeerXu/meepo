package transport_webrtc

import (
	"context"
)

func (t *WebrtcTransport) WaitResponse(ctx context.Context, in Message) (outs chan Message, err error) {
	logger := t.GetLogger().WithField("#method", "WaitResponse").WithFields(t.wrapMessage(in))
	resSess := t.parseResponseSession(in.Session)
	logger = logger.WithField("responseSession", resSess)
	lch, ok := t.polls.Load(resSess)
	if !ok {
		err = ErrSessionNotFoundFn(resSess)
		logger.WithError(err).Debugf("session not found")
		return nil, err
	}

	logger.Tracef("get outs")

	return lch.Ch, nil
}
