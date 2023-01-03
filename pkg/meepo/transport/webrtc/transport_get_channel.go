package transport_webrtc

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) GetChannel(ctx context.Context, id uint16) (meepo_interface.Channel, error) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":   "GetChannel",
		"channelID": id,
	})

	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	c, found := t.cs[id]
	if !found {
		err := ErrChannelNotFoundFn(id)
		logger.WithError(err).Debugf("channel not found")
		return nil, err
	}

	logger.Tracef("get channel")

	return c, nil
}
