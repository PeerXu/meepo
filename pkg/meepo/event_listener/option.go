package meepo_event_listener

import "github.com/PeerXu/meepo/pkg/lib/option"

const (
	DELIMITER     = "."
	WILDCARD      = "*"
	TAIL_WILDCARD = "*."

	OPTION_EVENT_LISTENER = "eventListener"
	OPTION_QUEUE_SIZE     = "queueSize"
)

var (
	WithEventListener, GetEventListener = option.New[EventListener](OPTION_EVENT_LISTENER)
)

type NewEventListenerOption = option.ApplyOption

func DefaultNewEventListenerOption() option.Option {
	return option.NewOption(map[string]any{
		OPTION_QUEUE_SIZE: 16,
	})
}

var (
	WithQueueSize, GetQueueSize = option.New[int](OPTION_QUEUE_SIZE)
)
