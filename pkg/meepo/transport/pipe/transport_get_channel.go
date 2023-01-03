package transport_pipe

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *PipeTransport) GetChannel(ctx context.Context, id uint16) (meepo_interface.Channel, error) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":   "GetChannel",
		"channelID": id,
	})

	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	c, found := t.cs[id]
	if !found {
		logger.Debugf("channel not found")
		return nil, transport_core.ErrChannelNotFoundFn(id)
	}

	logger.Tracef("get channel")

	return c, nil
}
