package teleportation_core

import (
	"net"

	"github.com/PeerXu/meepo/pkg/internal/listenerer"
	listenerer_interface "github.com/PeerXu/meepo/pkg/internal/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

type (
	Teleportation = meepo_interface.Teleportation
)

type teleportation struct {
	id         string
	mode       string
	logger     logging.Logger
	listener   listenerer_interface.Listener
	addr       addr.Addr
	sourceAddr net.Addr
	sinkAddr   net.Addr
	onAccept   OnTeleportationAcceptFunc
	onClose    OnTeleportationCloseFunc
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

	onCloseFn, err := GetOnTeleportationCloseFunc(o)
	if err != nil {
		return nil, err
	}

	tp := &teleportation{
		id:         id,
		mode:       mode,
		logger:     logger,
		addr:       addr,
		listener:   listener,
		sourceAddr: sourceAddr,
		sinkAddr:   sinkAddr,
		onAccept:   onAcceptFn,
		onClose:    onCloseFn,
	}

	go tp.acceptLoop()

	return tp, nil
}
