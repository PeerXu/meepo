package teleportation_core

import (
	"net"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/listenerer"
	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

type (
	Teleportation = meepo_interface.Teleportation
)

type teleportation struct {
	id                           string
	mode                         string
	logger                       logging.Logger
	listener                     listenerer_interface.Listener
	addr                         addr.Addr
	sourceAddr                   net.Addr
	sinkAddr                     net.Addr
	onAccept                     OnTeleportationAcceptFunc
	beforeCloseTeleportationHook BeforeCloseTeleportationHook
	afterCloseTeleportationHook  AfterCloseTeleportationHook
}

func NewTeleportation(opts ...NewTeleportationOption) (meepo_interface.Teleportation, error) {
	o := option.Apply(opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	id, err := well_known_option.GetID(o)
	if err != nil {
		return nil, err
	}

	mode, err := well_known_option.GetMode(o)
	if err != nil {
		return nil, err
	}

	listener, err := listenerer.GetListener(o)
	if err != nil {
		return nil, err
	}

	addr, err := well_known_option.GetAddr(o)
	if err != nil {
		return nil, err
	}

	sourceAddr, err := GetSourceAddr(o)
	if err != nil {
		return nil, err
	}

	sinkAddr, err := GetSinkAddr(o)
	if err != nil {
		return nil, err
	}

	onAcceptFn, err := GetOnTeleportationAcceptFunc(o)
	if err != nil {
		return nil, err
	}

	beforeNewTeleportationHook, _ := GetBeforeNewTeleportationHook(o)
	afterNewTeleportationHook, _ := GetAfterNewTeleportationHook(o)
	beforeCloseTeleportationHook, _ := GetBeforeCloseTeleportationHook(o)
	afterCloseTeleportationHook, _ := GetAfterCloseTeleportationHook(o)

	if h := beforeNewTeleportationHook; h != nil {
		if err = h(mode, sourceAddr.Network(), sourceAddr.String(), sinkAddr.Network(), sinkAddr.String()); err != nil {
			return nil, err
		}
	}

	tp := &teleportation{
		id:                           id,
		mode:                         mode,
		logger:                       logger,
		addr:                         addr,
		listener:                     listener,
		sourceAddr:                   sourceAddr,
		sinkAddr:                     sinkAddr,
		onAccept:                     onAcceptFn,
		beforeCloseTeleportationHook: beforeCloseTeleportationHook,
		afterCloseTeleportationHook:  afterCloseTeleportationHook,
	}

	go tp.acceptLoop()

	if h := afterNewTeleportationHook; h != nil {
		h(tp)
	}

	return tp, nil
}
