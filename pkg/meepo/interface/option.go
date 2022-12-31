package meepo_interface

import "github.com/PeerXu/meepo/pkg/internal/option"

type DialOption = option.ApplyOption

type CallOption = option.ApplyOption

type HandleOption = option.ApplyOption

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
	OPTION_MEEPO = "meepo"
)

var WithMeepo, GetMeepo = option.New[Meepo](OPTION_MEEPO)
