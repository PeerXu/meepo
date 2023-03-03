package transport_webrtc

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (t *WebrtcTransport) waitResponse(ctx context.Context, in Message) (outs chan Message, err error) {
	resSess := t.parseResponseSession(in.Session)
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":         "waitResponse",
		"responseSession": resSess,
	}).WithFields(t.wrapMessage(in))
	lch, ok := t.polls.Load(resSess)
	if !ok {
		err = ErrSessionNotFoundFn(resSess)
		logger.WithError(err).Debugf("session not found")
		return nil, err
	}
	return lch.Ch, nil
}
