package teleportation_core

import (
	"net"

	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	OPTION_SOURCE_ADDR                     = "sourceAddr"
	OPTION_SINK_ADDR                       = "sinkAddr"
	OPTION_ON_TELEPORTATION_ACCEPT_FUNC    = "onTeleportationAcceptFunc"
	OPTION_BEFORE_NEW_TELEPORTATION_HOOK   = "beforeNewTeleportationHook"
	OPTION_AFTER_NEW_TELEPORTATION_HOOK    = "afterNewTeleportationHook"
	OPTION_BEFORE_CLOSE_TELEPORTATION_HOOK = "beforeCloseTeleportationHook"
	OPTION_AFTER_CLOSE_TELEPORTATION_HOOK  = "afterCloseTeleportationHook"
)

type BeforeNewTeleportationHook = func(mode, sourceNetwork, sourceAddress, sinkNetwork, sinkAddress string, opts ...HookOption) error
type AfterNewTeleportationHook = func(Teleportation, ...HookOption)
type BeforeCloseTeleportationHook = func(Teleportation, ...HookOption) error
type AfterCloseTeleportationHook = func(Teleportation, ...HookOption)

type OnTeleportationAcceptFunc = func(Teleportation, listenerer_interface.Conn)

type NewTeleportationOption = option.ApplyOption

var (
	WithSourceAddr, GetSourceAddr                                     = option.New[net.Addr](OPTION_SOURCE_ADDR)
	WithSinkAddr, GetSinkAddr                                         = option.New[net.Addr](OPTION_SINK_ADDR)
	WithBeforeNewTeleportationHook, GetBeforeNewTeleportationHook     = option.New[BeforeNewTeleportationHook](OPTION_BEFORE_NEW_TELEPORTATION_HOOK)
	WithAfterNewTeleportationHook, GetAfterNewTeleportationHook       = option.New[AfterNewTeleportationHook](OPTION_AFTER_NEW_TELEPORTATION_HOOK)
	WithBeforeCloseTeleportationHook, GetBeforeCloseTeleportationHook = option.New[BeforeCloseTeleportationHook](OPTION_BEFORE_CLOSE_TELEPORTATION_HOOK)
	WithAfterCloseTeleportationHook, GetAfterCloseTeleportationHook   = option.New[AfterCloseTeleportationHook](OPTION_AFTER_CLOSE_TELEPORTATION_HOOK)
	WithOnTeleportationAcceptFunc, GetOnTeleportationAcceptFunc       = option.New[OnTeleportationAcceptFunc](OPTION_ON_TELEPORTATION_ACCEPT_FUNC)
)

type HookOption = option.ApplyOption
