package transport_pipe

import (
	"context"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *PipeTransport) ListChannels(ctx context.Context, opts ...meepo_interface.ListChannelsOption) ([]meepo_interface.Channel, error) {
	var cs []meepo_interface.Channel

	logger := t.GetLogger().WithField("#method", "ListChannels")

	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	for _, c := range t.cs {
		cs = append(cs, c)
	}

	logger.Tracef("list channels")

	return cs, nil
}
