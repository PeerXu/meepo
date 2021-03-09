package meepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/transport"
)

type Request interface {
	MessageGetter
}

type Response interface {
	MessageGetter
	Copier
}

type RequestHandler interface {
	Handle(dc transport.DataChannel, in interface{})
}

type MeepoRequestHandler func(transport.DataChannel, interface{})

func (h MeepoRequestHandler) Handle(dc transport.DataChannel, in interface{}) {
	h(dc, in)
}

type BroadcastRequestHandler interface {
	Handle(dc transport.DataChannel, in interface{})
	HandleBroadcast(dc transport.DataChannel, in interface{})
}

type MeepoBroadcastRequestHandler struct {
	handle          func(transport.DataChannel, interface{})
	handleBroadcast func(transport.DataChannel, interface{})
}

func (h *MeepoBroadcastRequestHandler) Handle(dc transport.DataChannel, in interface{}) {
	h.handle(dc, in)
}

func (h *MeepoBroadcastRequestHandler) HandleBroadcast(dc transport.DataChannel, in interface{}) {
	h.handleBroadcast(dc, in)
}

func (mp *Meepo) initHandlers() {
	mp.initRequestHandlers()
	mp.initBroadcastRequestHandlers()
}

func (mp *Meepo) registerRequestHandleFunc(name string, h func(transport.DataChannel, interface{})) {
	mp.registerRequestHandler(name, MeepoRequestHandler(h))
}

func (mp *Meepo) registerRequestHandler(name string, h RequestHandler) {
	mp.requestHandlersMtx.Lock()
	mp.requestHandlers[name] = h
	mp.requestHandlersMtx.Unlock()
}

func (mp *Meepo) getRequestHandler(name string) (RequestHandler, error) {
	mp.requestHandlersMtx.Lock()
	handler, ok := mp.requestHandlers[name]
	mp.requestHandlersMtx.Unlock()
	if !ok {
		return nil, UnsupportedRequestHandlerError
	}
	return handler, nil
}

func (mp *Meepo) registerBroadcastRequestHandler(name string, h BroadcastRequestHandler) {
	mp.broadcastRequestHandlersMtx.Lock()
	mp.broadcastRequestHandlers[name] = h
	mp.broadcastRequestHandlersMtx.Unlock()
}

func (mp *Meepo) getBroadcastRequestHandler(name string) (BroadcastRequestHandler, error) {
	mp.broadcastRequestHandlersMtx.Lock()
	handler, ok := mp.broadcastRequestHandlers[name]
	mp.broadcastRequestHandlersMtx.Unlock()
	if !ok {
		return nil, UnsupportedRequestHandlerError
	}
	return handler, nil
}

func (mp *Meepo) sysDataChannelLoop(dc transport.DataChannel) {
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "sysDataChannelLoop",
		"label":   dc.Label(),
	})

	defer logger.Tracef("done")

	logger.Tracef("started")
	buf := make([]byte, 65535)
	for {
		n, err := dc.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				logger.WithError(err).Debugf("failed to read data from data channel")
			}
			return
		}

		mp.onSysDataChannelMessage(dc, buf[:n])
	}
}

func (mp *Meepo) onSysDataChannelMessage(dc transport.DataChannel, buf []byte) {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onDataChannelMessageHandler",
		"label":   dc.Label(),
	})

	var in interface{}
	in, err = DecodeMessage(buf)
	if err != nil {
		logger.WithError(err).Debugf("failed to decode message")
		return
	}

	switch in.(MessageGetter).GetMessage().Type {
	case MESSAGE_TYPE_REQUEST:
		mp.handleRequest(dc, in)
	case MESSAGE_TYPE_BROADCAST_REQUEST:
		mp.handleBroadcastRequest(dc, in)
	case MESSAGE_TYPE_RESPONSE:
		if err = mp.handleResponse(dc, in); err != nil {
			logger.WithError(err).Debugf("failed to handle response")
		}
	}
}

func (mp *Meepo) recoverHandleRequest(m *Message) {
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "recoverHandleRequest",
		"method":  m.Method,
		"session": m.Session,
	})

	recovered := recover()
	if recovered == nil {
		return
	}

	switch recovered.(type) {
	case sendMessageError:
		err := recovered.(sendMessageError)
		logger.WithError(err).Warningf("send message error")
	case error:
		err := recovered.(error)
		logger.WithError(err).Debugf("failed to handle request")
	}
}

func (mp *Meepo) initRequestHandlers() {
	mp.registerRequestHandleFunc("ping", mp.onPing)
	mp.registerRequestHandleFunc("newTeleportation", mp.onNewTeleportation)
	mp.registerRequestHandleFunc("closeTeleportation", mp.onCloseTeleportation)
	mp.registerRequestHandleFunc("doTeleport", mp.onDoTeleport)
}

func (mp *Meepo) handleRequest(dc transport.DataChannel, in interface{}) {
	m := in.(MessageGetter).GetMessage()

	defer mp.recoverHandleRequest(m)

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "handleRequest",
		"method":  m.Method,
		"session": m.Session,
	})

	handler, err := mp.getRequestHandler(m.Method)
	if err != nil {
		logger.WithError(err).Debugf("failed to get request handler")
		return
	}

	handler.Handle(dc, in)
	logger.Tracef("done")
}

func (mp *Meepo) handleResponse(dc transport.DataChannel, in interface{}) error {
	m := in.(MessageGetter).GetMessage()

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "handleResponse",
		"method":  m.Method,
		"session": m.Session,
	})

	ch, unlock, err := mp.channelLocker.GetWithUnlock(m.Session)
	if err != nil {
		logger.WithError(err).Debugf("failed to get session channel")
		return err
	}
	defer unlock()

	ch <- in
	logger.Tracef("send response to session channel")

	return nil
}

func (mp *Meepo) sendRequest(id string, req interface{}) error {
	msg := req.(MessageGetter).GetMessage()

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "sendRequest",
		"id":      id,
		"session": msg.Session,
		"method":  msg.Method,
	})

	tp, err := mp.getTransport(id)
	// TODO(Peer): handle transport not found
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		return err
	}

	dc, err := mp.getOrCreateSysDataChannel(tp)
	if err != nil {
		logger.WithError(err).Debugf("failed to ensure sys data channel")
		return err
	}

	mp.sendMessage(dc, req)

	return nil
}

func (mp *Meepo) sendMessage(dc transport.DataChannel, msg interface{}) {
	buf, err := json.Marshal(msg)
	if err != nil {
		panic(SendMessageError(err))
	}

	if _, err = dc.Write(buf); err != nil {
		panic(SendMessageError(err))
	}
}

func (mp *Meepo) waitResponse(session int32) (interface{}, error) {
	ch, err := mp.channelLocker.Get(session)
	if err != nil {
		return nil, err
	}

	select {
	case res, ok := <-ch:
		if !ok {
			return nil, SessionChannelClosedError(session)
		}

		msgGetter, ok := res.(MessageGetter)
		if !ok {
			return nil, UnexpectedMessageError
		}

		msg := msgGetter.GetMessage()
		if msg.Error != "" {
			return nil, fmt.Errorf(msg.Error)
		}

		return res, nil
	case <-time.After(cast.ToDuration(mp.opt.Get("waitResponseTimeout").Inter())):
		return nil, WaitResponseTimeoutError
	}
}

func (mp *Meepo) getOrCreateSysDataChannel(tp transport.Transport) (transport.DataChannel, error) {
	var dc transport.DataChannel
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "getOrCreateSysDataChannel",
	})

	dc, err = tp.DataChannel("sys")
	if err != nil {
		if !errors.Is(err, transport.DataChannelNotFoundError) {
			logger.WithError(err).Debugf("failed to get data channel")
			return nil, err
		}

		var wg sync.WaitGroup
		wg.Add(1)
		if dc, err = tp.CreateDataChannel(
			"sys",
			transport.WithOrdered(true),
		); err != nil {
			logger.WithError(err).Debugf("failed to create data channel")
			return nil, err
		}
		dc.OnOpen(func() {
			go mp.sysDataChannelLoop(dc)
			wg.Done()
			logger.Tracef("sys data channel opened")
		})

		logger.Tracef("sys data channel created")

		wg.Wait()
	}

	return dc, nil
}

func (mp *Meepo) doRequest(id string, req interface{}) (interface{}, error) {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "doRequest",
		"peerID":  id,
	})

	msgGetter, ok := req.(MessageGetter)
	if !ok {
		logger.WithError(err).Debugf("failed to convert request to MessageGetter")
		return nil, UnexpectedMessageError
	}

	msg := msgGetter.GetMessage()
	logger = logger.WithFields(logrus.Fields{
		"method":  msg.Method,
		"session": msg.Session,
	})

	if err = mp.channelLocker.Acquire(msg.Session); err != nil {
		logger.WithError(err).Debugf("failed to acquire session channel")
		return nil, err
	}
	defer func() {
		mp.channelLocker.Release(msg.Session)
		logger.Tracef("release session channel")
	}()
	logger.Tracef("acquire session channel")

	if err = mp.sendRequest(id, req); err != nil {
		logger.WithError(err).Debugf("failed to send request")
		return nil, err
	}
	logger.Tracef("send request")

	res, err := mp.waitResponse(msg.Session)
	if err != nil {
		logger.WithError(err).Debugf("failed to wait response")
		return nil, err
	}
	logger.Tracef("receive response")

	logger.Tracef("done")

	return res, nil
}
