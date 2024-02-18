package meepo_core

import (
	"context"

	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/rand"
	meepo_debug_sdk "github.com/PeerXu/meepo/pkg/meepo/debug/sdk"
	meepo_debug_sdk_http "github.com/PeerXu/meepo/pkg/meepo/debug/sdk/http"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (mp *Meepo) enableMeepoDebugInterface(baseUrl string) {
	logger := mp.GetLogger().WithField("#method", "enableMeepoDebugInterface")

	mdi, err := meepo_debug_sdk.NewSDK("http", meepo_debug_sdk_http.WithBaseURL(baseUrl))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(mp.context())
	defer cancel()

	sess := rand.DefaultStringGenerator.Generate(8)
	evts, err := mp.WatchEvents(ctx, []string{"mpo.transport.state.*"}, meepo_interface.WithSession(sess))
	if err != nil {
		panic(err)
	}

	for evt := range evts {
		happenedAt := evt.HappenedAt
		target := cast.ToString(evt.Data["addr"])
		state := cast.ToString(evt.Data["state"])
		logger := logger.WithFields(logging.Fields{
			"happenedAt": happenedAt,
			"target":     target,
			"state":      state,
		})
		if err = mdi.TransportStateChange(
			mp.context(),
			happenedAt,
			mp.Addr().String(),
			target,
			state); err != nil {
			logger.Errorf("failed to send transport state change to meepo debug interface")
		} else {
			logger.Tracef("send transport state change to meepo debug interface")
		}
	}
}
