package meepo_core

import (
	"context"
	"errors"
	"io"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
	msync "github.com/PeerXu/meepo/pkg/lib/sync"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

const (
	WATCH_EVENTS_COMMAND_WATCH       = "watch"
	WATCH_EVENTS_COMMAND_UNWATCH     = "unwatch"
	WATCH_EVENTS_COMMAND_UNWATCH_ALL = "unwatchAll"
)

func (mp *Meepo) hdrStreamAPIWatchEvents(ctx context.Context, stm rpc_interface.Stream) error {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "hdrStreamAPIWatchEvents",
	})

	evts := make(chan meepo_interface.Event)
	defer close(evts)

	barrierEventsLoop := make(chan struct{})
	barrierCommandsLoop := make(chan struct{})
	go mp.hdrStreamAPIWatchEvents_eventsLoop(ctx, logger, evts, stm, barrierEventsLoop)
	go mp.hdrStreamAPIWatchEvents_commandsLoop(ctx, logger, evts, stm, barrierCommandsLoop)

	select {
	case <-barrierEventsLoop:
	case <-barrierCommandsLoop:
	}

	logger.Infof("watch events")

	return nil
}

func (mp *Meepo) hdrStreamAPIWatchEvents_eventsLoop(ctx context.Context, logger logging.Logger, evts chan meepo_interface.Event, stm rpc_interface.Stream, barrier chan struct{}) {
	defer close(barrier)

	logger = logger.WithFields(logging.Fields{
		"#method": "hdrStreamAPIWatchEvents_eventsLoop",
	})
	defer logger.Debugf("events loop done")

	logger.Debugf("events loop")
	for evt := range evts {
		logger := logger.WithFields(logging.Fields{
			"event":   evt.Name,
			"session": evt.Session,
		})

		msg, err := stm.Marshaler().Marshal(evt)
		if err != nil {
			logger.Errorf("failed to marshal event to message")
			break
		}

		if err = stm.SendMessage(msg); err != nil {
			logger.Errorf("failed to send event message")
			break
		}

		logger.Tracef("send event message")
	}
}

func (mp *Meepo) hdrStreamAPIWatchEvents_commandsLoop(ctx context.Context, logger logging.Logger, evts chan meepo_interface.Event, stm rpc_interface.Stream, barrier chan struct{}) {
	defer close(barrier)

	sessMap := msync.NewMap[string, context.CancelFunc]()

	logger = logger.WithFields(logging.Fields{
		"#method": "hdrStreamAPIWatchEvents_commandsLoop",
	})
	defer logger.Debugf("commands loop done")

	logger.Debugf("commands loop")
	for {
		msg, err := stm.RecvMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			logger.WithError(err).Errorf("failed to receive command message")
			return
		}

		var c sdk_interface.WatchEventsStream_Command
		if err = msg.Unmarshal(&c); err != nil {
			logger.WithError(err).Errorf("failed to unmarshal command message")
			return
		}
		logger = logger.WithField("command", c.Command)

		switch c.Command {
		case WATCH_EVENTS_COMMAND_WATCH:
			var c sdk_interface.WatchEventsStream_WatchCommand
			if err = msg.Unmarshal(&c); err != nil {
				logger.WithError(err).Errorf("failed to unmarshal watch command")
				return
			}

			logger := logger.WithFields(logging.Fields{
				"session":  c.Session,
				"policies": c.Policies,
			})

			_, found := sessMap.Load(c.Session)
			if found {
				logger.Warningf("watched session")
				continue
			}

			nctx, cancel := context.WithCancel(ctx)
			defer cancel()
			if _, err = mp.WatchEvents(nctx, c.Policies, meepo_interface.WithSession(c.Session), meepo_interface.WithEventChannel(evts)); err != nil {
				logger.WithError(err).Errorf("failed to watch events")
				return
			}

			sessMap.Store(c.Session, cancel)

			logger.Debugf("watch")
		case WATCH_EVENTS_COMMAND_UNWATCH:
			var c sdk_interface.WatchEventsStream_UnwatchCommand
			if err = msg.Unmarshal(&c); err != nil {
				logger.WithError(err).Errorf("failed to unmarshal unwatch command")
				return
			}

			logger := logger.WithFields(logging.Fields{
				"session": c.Session,
			})

			cancel, found := sessMap.LoadAndDelete(c.Session)
			if !found {
				logger.Warningf("session not found")
				continue
			}
			cancel()

			logger.Debugf("unwatch")
		case WATCH_EVENTS_COMMAND_UNWATCH_ALL:
			var sessions []string

			sessMap.Range(func(session string, cancel context.CancelFunc) bool {
				cancel()
				sessMap.Delete(session)
				sessions = append(sessions, session)
				return true
			})

			logger.WithField("sessions", sessions).Debugf("unwatch all")
		default:
			logger.Warningf("unsupported command")
		}
	}
}
