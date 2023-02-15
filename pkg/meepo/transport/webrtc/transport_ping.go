package transport_webrtc

import (
	"context"
	"time"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

type PingRequest struct {
	Session int64
}

type PingResponse struct {
	Session int64
}

func (t *WebrtcTransport) ping(ctx context.Context) (ttl time.Duration, err error) {
	var res PingResponse

	pingSess := t.randSrc.Int63()
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "ping",
		"session": pingSess,
	})
	pingAt := time.Now()

	if err = t.Call(ctx, SYS_METHOD_PING, &PingRequest{Session: pingSess}, &res, well_known_option.WithScope("sys")); err != nil {
		logger.WithError(err).Debugf("failed to ping")
		return time.Duration(0), err
	}

	ttl = time.Since(pingAt)
	logger = logger.WithField("ttl", ttl)

	if pingSess != res.Session {
		err = ErrInvalidPingSessionFn(pingSess, res.Session)
		logger.WithError(err).Debugf("invalid ping session")
		return time.Duration(0), err
	}

	logger.Tracef("ping")

	return
}

func (t *WebrtcTransport) onPing(ctx context.Context, _req any) (res any, err error) {
	req := _req.(*PingRequest)

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "onPing",
		"session": req.Session,
	})

	logger.Tracef("on ping")

	return &PingResponse{Session: req.Session}, nil
}
