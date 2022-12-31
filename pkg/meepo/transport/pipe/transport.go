package transport_pipe

import (
	"sync/atomic"

	"github.com/PeerXu/meepo/pkg/internal/dialer"
	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/option"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/lock"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type PipeTransport struct {
	addr meepo_interface.Addr

	state            atomic.Value
	currentChannelID uint32
	dialer           dialer.Dialer
	onClose          transport_core.OnTransportCloseFunc
	onReady          transport_core.OnTransportReadyFunc
	logger           logging.Logger

	csMtx lock.Locker
	cs    map[uint16]meepo_interface.Channel

	fnsMtx      lock.Locker
	fns         map[string]meepo_interface.HandleFunc
	marshaler   marshaler_interface.Marshaler
	unmarshaler marshaler_interface.Unmarshaler
}

func NewPipeTransport(opts ...meepo_interface.NewTransportOption) (meepo_interface.Transport, error) {
	o := option.Apply(opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	addr, err := well_known_option.GetAddr(o)
	if err != nil {
		return nil, err
	}

	dialer, err := dialer.GetDialer(o)
	if err != nil {
		return nil, err
	}

	onTransportClose, err := transport_core.GetOnTransportCloseFunc(o)
	if err != nil {
		return nil, err
	}

	onTransportReady, err := transport_core.GetOnTransportReadyFunc(o)
	if err != nil {
		return nil, err
	}

	mr, err := marshaler.GetMarshaler(o)
	if err != nil {
		return nil, err
	}

	umr, err := marshaler.GetUnmarshaler(o)
	if err != nil {
		return nil, err
	}

	t := &PipeTransport{
		addr:        addr,
		dialer:      dialer,
		onClose:     onTransportClose,
		onReady:     onTransportReady,
		logger:      logger,
		csMtx:       lock.NewLock(well_known_option.WithName("csMtx")),
		cs:          make(map[uint16]meepo_interface.Channel),
		fnsMtx:      lock.NewLock(well_known_option.WithName("fnsMtx")),
		fns:         make(map[string]meepo_interface.HandleFunc),
		marshaler:   mr,
		unmarshaler: umr,
	}
	t.setState(meepo_interface.TRANSPORT_STATE_NEW)
	defer onTransportReady(t) // nolint:errcheck
	go func() {
		t.setState(meepo_interface.TRANSPORT_STATE_NEW)
		t.setState(meepo_interface.TRANSPORT_STATE_CONNECTING)
		t.setState(meepo_interface.TRANSPORT_STATE_CONNECTED)
	}()

	return t, nil
}

func init() {
	transport_core.RegisterNewTransportFunc("pipe", NewPipeTransport)
}
