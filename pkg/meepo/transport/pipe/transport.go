package transport_pipe

import (
	"sync/atomic"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/lock"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

const (
	TRANSPORT_PIPE = "pipe"
)

type PipeTransport struct {
	transport_core.TransportHooks
	transport_core.ChannelHooks

	addr meepo_interface.Addr

	state            atomic.Value
	currentChannelID uint32
	dialer           dialer.Dialer
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

	if h, _ := transport_core.GetBeforeNewTransportHook(o); h != nil {
		if err = h(addr); err != nil {
			return nil, err
		}
	}

	dialer, err := dialer.GetDialer(o)
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
		onReady:     onTransportReady,
		logger:      logger,
		csMtx:       lock.NewLock(well_known_option.WithName("csMtx")),
		cs:          make(map[uint16]meepo_interface.Channel),
		fnsMtx:      lock.NewLock(well_known_option.WithName("fnsMtx")),
		fns:         make(map[string]meepo_interface.HandleFunc),
		marshaler:   mr,
		unmarshaler: umr,
	}

	transport_core.ApplyTransportHooks(o, &t.TransportHooks)
	transport_core.ApplyChannelHooks(o, &t.ChannelHooks)

	t.setState(meepo_interface.TRANSPORT_STATE_NEW)
	go func() {
		defer onTransportReady(t) // nolint:errcheck
		t.setState(meepo_interface.TRANSPORT_STATE_NEW)
		t.setState(meepo_interface.TRANSPORT_STATE_CONNECTING)
		t.setState(meepo_interface.TRANSPORT_STATE_CONNECTED)
	}()

	if h := t.AfterNewTransportHook; h != nil {
		h(t)
	}

	return t, nil
}

func init() {
	transport_core.RegisterNewTransportFunc(TRANSPORT_PIPE, NewPipeTransport)
}
