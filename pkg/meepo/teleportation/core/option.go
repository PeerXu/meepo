package teleportation_core

import (
	"net"

	listenerer_interface "github.com/PeerXu/meepo/pkg/internal/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/internal/option"
)

const (
	OPTION_SOURCE_ADDR                  = "sourceAddr"
	OPTION_SINK_ADDR                    = "sinkAddr"
	OPTION_ON_TELEPORTATION_ACCEPT_FUNC = "onTeleportationAcceptFunc"
	OPTION_ON_TELEPORTATION_CLOSE_FUNC  = "onTeleportationCloseFunc"
)

type OnTeleportationAcceptFunc = func(Teleportation, listenerer_interface.Conn)
type OnTeleportationCloseFunc = func(Teleportation) error

type NewTeleportationOption = option.ApplyOption

var (
	WithSourceAddr, GetSourceAddr                               = option.New[net.Addr](OPTION_SOURCE_ADDR)
	WithSinkAddr, GetSinkAddr                                   = option.New[net.Addr](OPTION_SINK_ADDR)
	WithOnTeleportationAcceptFunc, GetOnTeleportationAcceptFunc = option.New[OnTeleportationAcceptFunc](OPTION_ON_TELEPORTATION_ACCEPT_FUNC)
	WithOnTeleportationCloseFunc, GetOnTeleportationCloseFunc   = option.New[OnTeleportationCloseFunc](OPTION_ON_TELEPORTATION_CLOSE_FUNC)
)
