package meepo_interface

import "github.com/PeerXu/meepo/pkg/lib/option"

type DialOption = option.ApplyOption

type CallOption = option.ApplyOption

type HandleOption = option.ApplyOption

type WatchEventsOption = option.ApplyOption

type NewTransportOption = option.ApplyOption

type ListTransportsOption = option.ApplyOption

type GetTransportOption = option.ApplyOption

type NewChannelOption = option.ApplyOption

type ListChannelsOption = option.ApplyOption

type GetChannelOption = option.ApplyOption

type NewTeleportationOption = option.ApplyOption

type ListTeleportationsOption = option.ApplyOption

type GetTeleportationOption = option.ApplyOption

type TeleportOption = option.ApplyOption

const (
	OPTION_MEEPO         = "meepo"
	OPTION_SESSION       = "session"
	OPTION_EVENT_CHANNEL = "eventChannel"
)

var (
	WithMeepo, GetMeepo               = option.New[Meepo](OPTION_MEEPO)
	WithSession, GetSession           = option.New[string](OPTION_SESSION)
	WithEventChannel, GetEventChannel = option.New[chan Event](OPTION_EVENT_CHANNEL)
)
