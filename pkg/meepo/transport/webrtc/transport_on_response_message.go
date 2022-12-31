package transport_webrtc

import "strconv"

func (t *WebrtcTransport) onResponseMessage(m Message) {
	logger := t.GetLogger().WithField("#method", "onResponseMessage").WithFields(t.wrapMessage(m))
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("catch error")
			panic(r)
		}
	}()

	lch, ok := t.polls.Load(m.Session)
	if !ok {
		logger.Debugf("session not in polls")
		return
	}

	if lch.Ch != nil {
		lch.Ch <- m
	}

	logger.Tracef("on response message")
}

func (t *WebrtcTransport) isResponseSession(sess string) bool {
	sessU64, _ := strconv.ParseUint(sess, 16, 32)
	return sessU64%2 == 0
}
