package webrtc_transport

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/transport"
)

type AnswerHook func(*webrtc.SessionDescription, error)
type OfferHook func(*webrtc.SessionDescription) (*webrtc.SessionDescription, error)

func newNewWebrtcTransportOptions() objx.Map {
	return objx.New(map[string]interface{}{
		"gatherTimeout": 17 * time.Second,
	})
}

func newCreateDataChannelOptions() objx.Map {
	return objx.New(map[string]interface{}{})
}

type WebrtcTransport struct {
	opt       objx.Map
	logger    logrus.FieldLogger
	peerID    string
	handleIdx transport.HandleID

	pc *webrtc.PeerConnection

	channels map[string]transport.DataChannel

	onTransportStateChangeHandler func(transport.TransportState)
	onDataChannelCreateHandlers   map[string]transport.OnDataChannelCreateHandler

	onTransportStateHandlers map[transport.TransportState]map[transport.HandleID]transport.OnTransportStateHandler

	channelsMtx                      sync.Mutex
	onTransportStateChangeHandlerMtx sync.Mutex
	onDataChannelCreateHandlersMtx   sync.Mutex
	onTransportStateHandlersMtx      sync.Mutex

	err     error
	errOnce sync.Once
}

func (wt *WebrtcTransport) PeerID() string {
	if wt.peerID == "" {
		wt.peerID = cast.ToString(wt.opt.Get("peerID").Inter())
	}

	return wt.peerID
}

func (wt *WebrtcTransport) getRawLogger() logrus.FieldLogger {
	return wt.logger
}

func (wt *WebrtcTransport) getLogger() logrus.FieldLogger {
	return wt.logger.WithFields(logrus.Fields{
		"#instance": "WebrtcTransport",
		"peerID":    wt.PeerID(),
	})
}

func (wt *WebrtcTransport) setError(err error) {
	wt.errOnce.Do(func() { wt.err = err })
}

func (wt *WebrtcTransport) setDataChannel(label string, dc transport.DataChannel) {
	wt.channelsMtx.Lock()
	wt.channels[label] = dc
	wt.channelsMtx.Unlock()
}

func (wt *WebrtcTransport) Err() error {
	return wt.err
}

func (wt *WebrtcTransport) Close() error {
	return wt.pc.Close()
}

func (wt *WebrtcTransport) OnTransportStateChange(f func(transport.TransportState)) {
	wt.onTransportStateChangeHandlerMtx.Lock()
	wt.onTransportStateChangeHandler = f
	wt.onTransportStateChangeHandlerMtx.Unlock()
}

func (wt *WebrtcTransport) OnTransportState(s transport.TransportState, f func(transport.HandleID)) transport.HandleID {
	hid := atomic.AddUint32(&wt.handleIdx, 1)

	wt.onTransportStateHandlersMtx.Lock()
	hm, ok := wt.onTransportStateHandlers[s]
	if !ok {
		hm = make(map[transport.HandleID]transport.OnTransportStateHandler)
		wt.onTransportStateHandlers[s] = hm
	}
	hm[hid] = f
	wt.onTransportStateHandlersMtx.Unlock()

	if wt.TransportState() == s {
		go f(hid)
	}

	return hid
}

func (wt *WebrtcTransport) UnsetOnTransportState(s transport.TransportState, hid transport.HandleID) {
	wt.onTransportStateHandlersMtx.Lock()
	if hm, ok := wt.onTransportStateHandlers[s]; ok {
		delete(hm, hid)
	}
	wt.onTransportStateHandlersMtx.Unlock()
}

func (wt *WebrtcTransport) onTransportState(s transport.TransportState) {
	logger := wt.getLogger().WithFields(logrus.Fields{
		"#method": "onTransportState",
		"peerID":  wt.PeerID(),
		"state":   s.String(),
	})

	wt.onTransportStateHandlersMtx.Lock()
	hm, ok := wt.onTransportStateHandlers[s]
	if ok {
		for hid, f := range hm {
			go func(hid transport.HandleID, f func(transport.HandleID)) {
				f(hid)
				logger.WithField("handleID", hid).Tracef("handle on transport state")
			}(hid, f)
		}
	}
	wt.onTransportStateHandlersMtx.Unlock()
}

func (wt *WebrtcTransport) TransportState() transport.TransportState {
	if wt.pc == nil {
		return transport.TransportStateNew
	}

	return transport.TransportState(wt.pc.ConnectionState())
}

func (wt *WebrtcTransport) DataChannels() ([]transport.DataChannel, error) {
	wt.channelsMtx.Lock()
	defer wt.channelsMtx.Unlock()

	var dcs []transport.DataChannel
	for _, dc := range wt.channels {
		dcs = append(dcs, dc)
	}

	return dcs, nil
}

func (wt *WebrtcTransport) DataChannel(label string) (transport.DataChannel, error) {
	wt.channelsMtx.Lock()
	defer wt.channelsMtx.Unlock()

	dc := wt.channels[label]
	if dc == nil {
		return nil, transport.DataChannelNotFoundError
	}

	return dc, nil
}

func (wt *WebrtcTransport) CreateDataChannel(label string, opts ...transport.CreateDataChannelOption) (transport.DataChannel, error) {
	var ordered bool

	logger := wt.getLogger().WithFields(logrus.Fields{
		"#method": "CreateDataChannel",
		"label":   label,
	})

	o := newCreateDataChannelOptions()

	for _, opt := range opts {
		opt(o)
	}

	ordered = cast.ToBool(o.Get("ordered").Inter())

	dc, err := wt.pc.CreateDataChannel(label, &webrtc.DataChannelInit{
		Ordered: &ordered,
	})
	if err != nil {
		logger.WithError(err).Debugf("failed to create data channel")
		return nil, err
	}
	wdc := NewWebrtcDataChannel(wt.getRawLogger(), dc, wt)
	logger.Tracef("create PeerConnection.DataChannel")

	wt.setDataChannel(label, wdc)
	logger.Debugf("data channel created")

	return wdc, nil
}

func (wt *WebrtcTransport) OnDataChannelCreate(label string, h func(transport.DataChannel)) {
	wt.onDataChannelCreateHandlersMtx.Lock()
	wt.onDataChannelCreateHandlers[label] = h
	wt.onDataChannelCreateHandlersMtx.Unlock()
}

func (wt *WebrtcTransport) initAsAnswerer() {
	var err error
	var closer func() error

	defer func() {
		if err != nil {
			wt.setError(err)
		}
	}()

	logger := wt.getLogger().WithField("#method", "initAsAnswerer")

	logger.Tracef("start")

	api := wt.opt.Get("webrtcAPI").Inter().(*webrtc.API)
	iceServers := wt.opt.Get("iceServers").Inter().([]string)
	offer := wt.opt.Get("offer").Inter().(*webrtc.SessionDescription)
	answerHook := wt.opt.Get("answerHook").Inter().(AnswerHook)

	wt.pc, err = api.NewPeerConnection(webrtc.Configuration{
		ICEServers: unmarshalICEServers(iceServers),
	})
	if err != nil {
		logger.WithError(err).Debugf("failed to new peer connection")
		return
	}
	closer = wt.pc.Close
	defer func() {
		if closer != nil {
			if err = closer(); err != nil {
				logger.WithError(err).Debugf("failed to close peer connection")
			}
		}
	}()
	wt.pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		label := dc.Label()

		innerLogger := wt.getLogger().WithFields(logrus.Fields{
			"#method": "OnDataChannel",
			"as":      "answerer",
			"label":   label,
		})

		wdc := NewWebrtcDataChannel(wt.getRawLogger(), dc, wt)
		innerLogger.Tracef("create PeerConnection.DataChannel")

		wt.setDataChannel(label, wdc)

		wt.onDataChannelCreateHandlersMtx.Lock()
		handler := wt.onDataChannelCreateHandlers[label]
		wt.onDataChannelCreateHandlersMtx.Unlock()

		if handler != nil {
			handler(wdc)
			innerLogger.Tracef("handle data channel created")
		}

		innerLogger.Tracef("data channel created")
	})
	wt.pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		innerLogger := wt.getLogger().WithFields(logrus.Fields{
			"#method": "OnConnectionStateChange",
			"as":      "answerer",
			"state":   s,
		})

		wt.onDataChannelCreateHandlersMtx.Lock()
		handler := wt.onTransportStateChangeHandler
		wt.onDataChannelCreateHandlersMtx.Unlock()
		if handler != nil {
			handler(transport.TransportState(s))
			innerLogger.Tracef("handle peer connection state changed")
		}

		wt.onTransportState(transport.TransportState(s))

		innerLogger.Tracef("peer connection state changed")
	})
	logger.Tracef("new peer connection")

	if err = wt.pc.SetRemoteDescription(*offer); err != nil {
		logger.WithError(err).Debugf("failed to set remote description")
		return
	}

	answer, err := wt.pc.CreateAnswer(nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create answer")
		return
	}
	logger.Tracef("create answer")

	gatherFinished := webrtc.GatheringCompletePromise(wt.pc)

	err = wt.pc.SetLocalDescription(answer)
	if err != nil {
		logger.WithError(err).Debugf("failed to set local description")
		return
	}
	logger.Tracef("set local description")

	logger.Tracef("gather start")
	select {
	case <-gatherFinished:
		logger.Tracef("gather done")
	case <-time.After(cast.ToDuration(wt.opt.Get("gatherTimeout").Inter())):
		err = GatherTimeoutError
		logger.WithError(err).Debugf("gather timeout")
		answerHook(nil, err)
		return
	}

	logger.Tracef("answer hook start")
	answerHook(wt.pc.LocalDescription(), nil)
	logger.Tracef("answer hook done")

	closer = nil

	logger.Tracef("done")
}

func (wt *WebrtcTransport) initAsOfferer() {
	var err error
	var closer func() error

	defer func() {
		if err != nil {
			wt.setError(err)
		}
	}()

	logger := wt.getLogger().WithField("#method", "initAsOfferer")

	logger.Tracef("start")

	api := wt.opt.Get("webrtcAPI").Inter().(*webrtc.API)
	iceServers := wt.opt.Get("iceServers").Inter().([]string)
	offerHook := wt.opt.Get("offerHook").Inter().(OfferHook)

	wt.pc, err = api.NewPeerConnection(webrtc.Configuration{
		ICEServers: unmarshalICEServers(iceServers),
	})
	if err != nil {
		logger.WithError(err).Debugf("failed to new peer connection")
		return
	}
	closer = wt.pc.Close
	defer func() {
		if closer != nil {
			if err = closer(); err != nil {
				logger.WithError(err).Debugf("failed to close peer connection")
			}
		}
	}()

	wt.pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		label := dc.Label()

		innerLogger := wt.getLogger().WithFields(logrus.Fields{
			"#method": "OnDataChannel",
			"as":      "offerer",
			"label":   label,
		})

		wdc := NewWebrtcDataChannel(wt.getRawLogger(), dc, wt)
		wt.setDataChannel(label, wdc)

		wt.onDataChannelCreateHandlersMtx.Lock()
		handler := wt.onDataChannelCreateHandlers[label]
		wt.onDataChannelCreateHandlersMtx.Unlock()

		if handler != nil {
			handler(wdc)
			innerLogger.Tracef("handle data channel created")
		}

		innerLogger.Tracef("data channel opened")
	})
	wt.pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		innerLogger := wt.getLogger().WithFields(logrus.Fields{
			"#method": "OnConnectionStateChange",
			"as":      "offerer",
			"state":   s,
		})

		wt.onTransportStateChangeHandlerMtx.Lock()
		handler := wt.onTransportStateChangeHandler
		wt.onTransportStateChangeHandlerMtx.Unlock()

		if handler != nil {
			handler(transport.TransportState(s))
			innerLogger.Tracef("handle data channel state changed")
		}

		wt.onTransportState(transport.TransportState(s))

		innerLogger.Tracef("data channel state changed")
	})
	logger.Tracef("new peer connection")

	// HACK(Peer): pion need at less one data channel
	dc, err := wt.pc.CreateDataChannel("_IGNORE_", nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create data channel")
		return
	}
	defer dc.Close()
	logger.Tracef("create _IGNORE_ data channel")

	offer, err := wt.pc.CreateOffer(nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create offer")
		return
	}
	logger.Tracef("create offer")

	gatherFinished := webrtc.GatheringCompletePromise(wt.pc)

	if err = wt.pc.SetLocalDescription(offer); err != nil {
		logger.WithError(err).Debugf("failed to set local description")
		return
	}
	logger.Tracef("set local description")

	logger.Tracef("gather start")
	select {
	case <-gatherFinished:
		logger.Tracef("gather done")
	case <-time.After(cast.ToDuration(wt.opt.Get("gatherTimeout").Inter())):
		err = GatherTimeoutError
		logger.WithError(err).Debugf("gather timeout")
		return
	}

	logger.Tracef("offer hook start")
	answer, err := offerHook(wt.pc.LocalDescription())
	if err != nil {
		logger.WithError(err).Debugf("failed to offer hook")
		return
	}
	logger.Tracef("offer hook done")

	if err = wt.pc.SetRemoteDescription(*answer); err != nil {
		logger.WithError(err).Debugf("failed to set remote description")
		return
	}
	logger.Tracef("set remote description")

	closer = nil

	logger.Tracef("done")
}

func NewWebrtcTransport(opts ...transport.NewTransportOption) (transport.Transport, error) {
	var ok bool
	var logger logrus.FieldLogger
	var role string

	o := newNewWebrtcTransportOptions()

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

	if _, ok = o.Get("webrtcAPI").Inter().(*webrtc.API); !ok {
		return nil, fmt.Errorf("Require webrtcAPI")
	}

	if _, ok = o.Get("iceServers").Inter().([]string); !ok {
		return nil, fmt.Errorf("Require iceServers")
	}

	if role, ok = o.Get("role").Inter().(string); !ok {
		return nil, fmt.Errorf("Require role")
	}

	wt := &WebrtcTransport{
		opt:    o,
		logger: logger,

		channels: make(map[string]transport.DataChannel),

		onDataChannelCreateHandlers: make(map[string]transport.OnDataChannelCreateHandler),
		onTransportStateHandlers:    make(map[transport.TransportState]map[transport.HandleID]transport.OnTransportStateHandler),
	}

	switch role {
	case "answerer":
		if _, ok = o.Get("offer").Inter().(*webrtc.SessionDescription); !ok {
			return nil, fmt.Errorf("Require offer")
		}

		if _, ok = o.Get("answerHook").Inter().(AnswerHook); !ok {
			return nil, fmt.Errorf("Require answerHook")
		}

		go wt.initAsAnswerer()
	case "offerer":
		if _, ok = o.Get("offerHook").Inter().(OfferHook); !ok {
			return nil, fmt.Errorf("Require offerHook")
		}

		go wt.initAsOfferer()
	default:
		return nil, UnsupportedRoleError(role)
	}

	return wt, nil
}

func init() {
	transport.RegisterNewTransportFunc("webrtc", NewWebrtcTransport)
}
