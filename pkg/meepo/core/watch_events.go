package meepo_core

import (
	"context"
	"time"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	meepo_eventloop_interface "github.com/PeerXu/meepo/pkg/meepo/eventloop/interface"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (mp *Meepo) WatchEvents(ctx context.Context, policies []string, opts ...meepo_interface.WatchEventsOption) (chan meepo_interface.Event, error) {
	o := option.ApplyWithDefault(mp.defaultWatchEventsOptions(), opts...)
	sess, _ := meepo_interface.GetSession(o)
	evtCh, _ := meepo_interface.GetEventChannel(o)
	withEventChannel := evtCh != nil

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":          "WatchEvents",
		"policies":         policies,
		"session":          sess,
		"withEventChannel": withEventChannel,
	})

	if !withEventChannel {
		evtCh = make(chan meepo_interface.Event)
	}
	cbids := []string{}
	cb := func(e meepo_eventloop_interface.Event) {
		logger.WithField("name", e.Name()).Tracef("receive event")
		evtCh <- meepo_interface.Event{
			Session:    sess,
			Name:       e.Name(),
			ID:         e.ID(),
			HappenedAt: e.HappenedAt().Format(time.RFC3339Nano),
			Data:       e.Data(),
		}
	}
	for _, p := range policies {
		cbid := mp.eventListener.Listen(p, cb)
		cbids = append(cbids, cbid)
	}
	go func() {
		<-ctx.Done()
		for _, cbid := range cbids {
			mp.eventListener.Unlisten(cbid)
		}
		if !withEventChannel {
			close(evtCh)
		}

		logger.Tracef("watch events done")
	}()

	logger.Tracef("watch events")
	return evtCh, nil
}
