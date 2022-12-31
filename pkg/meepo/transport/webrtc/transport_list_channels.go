package transport_webrtc

import (
	"context"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) ListChannels(ctx context.Context, opts ...meepo_interface.ListChannelsOption) (cs []meepo_interface.Channel, err error) {
	logger := t.GetLogger().WithField("#method", "ListChannels")

	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	for _, c := range t.cs {
		cs = append(cs, c)
	}

	logger.Tracef("list channels")

	return
}
