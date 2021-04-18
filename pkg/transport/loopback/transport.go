package loopback_transport

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/transport"
)

func newNewLoopbackTransportOptions() objx.Map {
	return objx.New(map[string]interface{}{})
}

type LoopbackTransport struct {
	opt       objx.Map
	logger    logrus.FieldLogger
	peerID    string
	handleIdx transport.HandleID

	state    transport.TransportState
	stateMtx sync.Locker

	channels    map[string]*LoopbackDataChannel
	channelsMtx sync.Locker

	onTransportStateChangeHandler    func(transport.TransportState)
	onTransportStateChangeHandlerMtx sync.Locker

	onDataChannelCreateHandlers    map[string]transport.OnDataChannelCreateHandler
	onDataChannelCreateHandlersMtx sync.Locker

	onTransportStateHandlers    map[transport.TransportState]map[transport.HandleID]transport.OnTransportStateHandler
	onTransportStateHandlersMtx sync.Locker

	err     error
	errOnce sync.Once
}

func NewLoopbackTransport(opts ...transport.NewTransportOption) (transport.Transport, error) {
	var ok bool
	var logger logrus.FieldLogger

	o := newNewLoopbackTransportOptions()

	for _, opt := range opts {
		opt(o)
	}

	if logger, ok = o.Get("logger").Inter().(logrus.FieldLogger); !ok {
		logger = logrus.New()
	}

	if _, ok = o.Get("id").Inter().(string); !ok {
		return nil, fmt.Errorf("Require id")
	}

	if _, ok = o.Get("peerID").Inter().(string); !ok {
		return nil, fmt.Errorf("Require peerID")
	}

	lt := &LoopbackTransport{
		opt:    o,
		logger: logger,

		state:    transport.TransportStateNew,
		stateMtx: new(sync.Mutex),

		channels:    make(map[string]*LoopbackDataChannel),
		channelsMtx: new(sync.Mutex),

		onTransportStateChangeHandlerMtx: new(sync.Mutex),

		onDataChannelCreateHandlers:    make(map[string]transport.OnDataChannelCreateHandler),
		onDataChannelCreateHandlersMtx: new(sync.Mutex),

		onTransportStateHandlers:    make(map[transport.TransportState]map[transport.HandleID]transport.OnTransportStateHandler),
		onTransportStateHandlersMtx: new(sync.Mutex),
	}

	go lt.onTransportCreate()

	return lt, nil
}

func (lt *LoopbackTransport) PeerID() string {
	if lt.peerID == "" {
		lt.peerID = cast.ToString(lt.opt.Get("peerID").Inter())
	}

	return lt.peerID
}

func (lt *LoopbackTransport) getRawLogger() logrus.FieldLogger {
	return lt.logger
}

func (lt *LoopbackTransport) getLogger() logrus.FieldLogger {
	return lt.logger.WithFields(logrus.Fields{
		"#instance": "LoopbackTransport",
		"peerID":    lt.PeerID(),
	})
}

func (lt *LoopbackTransport) Err() error {
	return lt.err
}

func (lt *LoopbackTransport) isClosed() bool {
	lt.stateMtx.Lock()
	defer lt.stateMtx.Unlock()
	return lt.isClosedNL()
}

func (lt *LoopbackTransport) isClosedNL() bool {
	return lt.state == transport.TransportStateClosed
}

func (lt *LoopbackTransport) Close() error {
	lt.stateMtx.Lock()
	if lt.isClosedNL() {
		return nil
	}
	lt.state = transport.TransportStateClosed
	lt.stateMtx.Unlock()

	go func() {
		lt.channelsMtx.Lock()
		for _, ch := range lt.channels {
			ch.Close()
		}
		lt.channelsMtx.Unlock()

		lt.onTransportStateChange(transport.TransportStateClosed)
		lt.onTransportState(transport.TransportStateClosed)
	}()

	return nil
}

func (lt *LoopbackTransport) onTransportStateChange(s transport.TransportState) {
	logger := lt.getLogger().WithFields(logrus.Fields{
		"#method": "onTransportStateChange",
	})

	lt.onTransportStateChangeHandlerMtx.Lock()
	handler := lt.onTransportStateChangeHandler
	lt.onTransportStateChangeHandlerMtx.Unlock()

	if handler != nil {
		handler(s)
		logger.Tracef("handle loopback state changed")
	}
}

func (lt *LoopbackTransport) onTransportState(s transport.TransportState) {
	logger := lt.getLogger().WithFields(logrus.Fields{
		"#method": "onTransportState",
		"peerID":  lt.PeerID(),
		"state":   s.String(),
	})

	lt.onTransportStateHandlersMtx.Lock()
	hm, ok := lt.onTransportStateHandlers[s]
	if ok {
		for hid, f := range hm {
			go func(hid transport.HandleID, f func(transport.HandleID)) {
				f(hid)
				logger.WithField("handleID", hid).Tracef("handle on transport state")
			}(hid, f)
		}
	}
	lt.onTransportStateHandlersMtx.Unlock()
}

func (lt *LoopbackTransport) OnTransportStateChange(f func(transport.TransportState)) {
	lt.onTransportStateChangeHandlerMtx.Lock()
	lt.onTransportStateChangeHandler = f
	lt.onTransportStateChangeHandlerMtx.Unlock()
}

func (lt *LoopbackTransport) OnTransportState(s transport.TransportState, f func(hid transport.HandleID)) transport.HandleID {
	hid := atomic.AddUint32(&lt.handleIdx, 1)

	lt.onTransportStateHandlersMtx.Lock()
	hm, ok := lt.onTransportStateHandlers[s]
	if !ok {
		hm = make(map[transport.HandleID]transport.OnTransportStateHandler)
		lt.onTransportStateHandlers[s] = hm
	}
	hm[hid] = f
	lt.onTransportStateHandlersMtx.Unlock()

	if lt.TransportState() == s {
		go f(hid)
	}

	return hid
}

func (lt *LoopbackTransport) UnsetOnTransportState(s transport.TransportState, hid transport.HandleID) {
	lt.onTransportStateHandlersMtx.Lock()
	if hm, ok := lt.onTransportStateHandlers[s]; ok {
		delete(hm, hid)
	}
	lt.onTransportStateHandlersMtx.Unlock()
}

func (lt *LoopbackTransport) TransportState() transport.TransportState {
	return lt.state
}

func (lt *LoopbackTransport) DataChannels() ([]transport.DataChannel, error) {
	lt.channelsMtx.Lock()
	defer lt.channelsMtx.Unlock()

	var dcs []transport.DataChannel
	for _, dc := range lt.channels {
		dcs = append(dcs, dc)
	}

	return dcs, nil
}

func (lt *LoopbackTransport) DataChannel(label string) (transport.DataChannel, error) {
	lt.channelsMtx.Lock()
	defer lt.channelsMtx.Unlock()

	dc := lt.channels[label]
	if dc == nil {
		return nil, transport.DataChannelNotFoundError
	}

	return dc.Left(), nil
}

func (lt *LoopbackTransport) setDataChannel(label string, dc *LoopbackDataChannel) {
	lt.channelsMtx.Lock()
	lt.channels[label] = dc
	lt.channelsMtx.Unlock()
}

func (lt *LoopbackTransport) CreateDataChannel(label string, opts ...transport.CreateDataChannelOption) (transport.DataChannel, error) {
	logger := lt.getLogger().WithFields(logrus.Fields{
		"#method": "CreateDataChannel",
		"label":   label,
	})

	dc := NewLoopbackDataChannel(label, lt, lt.getRawLogger())
	logger.Tracef("create LoopbackDataChannel")

	lt.setDataChannel(label, dc)
	logger.Debugf("data channel created")

	go lt.onDataChannel(dc)

	return dc.Left(), nil
}

func (lt *LoopbackTransport) handleDataChannelCreate(dc transport.DataChannel) (done chan struct{}) {
	lt.onDataChannelCreateHandlersMtx.Lock()
	handler := lt.onDataChannelCreateHandlers[dc.Label()]
	lt.onDataChannelCreateHandlersMtx.Unlock()

	done = make(chan struct{})
	if handler == nil {
		close(done)
		return
	}

	go func() {
		handler(dc)
		close(done)
	}()

	return
}

func (lt *LoopbackTransport) onDataChannel(dc *LoopbackDataChannel) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-lt.handleDataChannelCreate(dc.Right())
		wg.Done()
	}()
	wg.Wait()

	dc.setState(transport.DataChannelStateOpen)
	go dc.Left().onOpen()
	go dc.Right().onOpen()
}

func (lt *LoopbackTransport) OnDataChannelCreate(label string, f func(transport.DataChannel)) {
	lt.onDataChannelCreateHandlersMtx.Lock()
	lt.onDataChannelCreateHandlers[label] = f
	lt.onDataChannelCreateHandlersMtx.Unlock()
}

func (lt *LoopbackTransport) handleTransportCreate() (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		lt.stateMtx.Lock()
		lt.state = transport.TransportStateConnected
		lt.stateMtx.Unlock()

		lt.onTransportStateChange(transport.TransportStateConnected)
		lt.onTransportState(transport.TransportStateConnected)

		close(done)
	}()
	return
}

func (lt *LoopbackTransport) onTransportCreate() {
	<-lt.handleTransportCreate()
}

func init() {
	transport.RegisterNewTransportFunc("loopback", NewLoopbackTransport)
}
