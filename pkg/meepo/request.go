package meepo

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/transport"
)

type RequestHandler interface {
	Handle(transport.DataChannel, packet.Packet)
}

type MeepoRequestHandler func(transport.DataChannel, packet.Packet)

func (h MeepoRequestHandler) Handle(dc transport.DataChannel, in packet.Packet) {
	h(dc, in)
}

type BroadcastRequestHandler interface {
	Handle(transport.DataChannel, packet.BroadcastPacket)
	HandleBroadcast(transport.DataChannel, packet.BroadcastPacket)
}

type MeepoBroadcastRequestHandler struct {
	handle          func(transport.DataChannel, packet.BroadcastPacket)
	handleBroadcast func(transport.DataChannel, packet.BroadcastPacket)
}

func (h *MeepoBroadcastRequestHandler) Handle(dc transport.DataChannel, in packet.BroadcastPacket) {
	h.handle(dc, in)
}

func (h *MeepoBroadcastRequestHandler) HandleBroadcast(dc transport.DataChannel, in packet.BroadcastPacket) {
	h.handleBroadcast(dc, in)
}

func (mp *Meepo) initHandlers() {
	mp.initRequestHandlers()
	mp.initBroadcastRequestHandlers()
}

func (mp *Meepo) registerRequestHandleFunc(m Method, h func(transport.DataChannel, packet.Packet)) {
	mp.registerRequestHandler(m, MeepoRequestHandler(h))
}

func (mp *Meepo) registerRequestHandler(m Method, h RequestHandler) {
	mp.requestHandlersMtx.Lock()
	mp.requestHandlers[m] = h
	mp.requestHandlersMtx.Unlock()
}

func (mp *Meepo) getRequestHandler(m Method) (RequestHandler, error) {
	mp.requestHandlersMtx.Lock()
	handler, ok := mp.requestHandlers[m]
	mp.requestHandlersMtx.Unlock()
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	return handler, nil
}

func (mp *Meepo) registerBroadcastRequestHandler(m Method, h BroadcastRequestHandler) {
	mp.broadcastRequestHandlersMtx.Lock()
	mp.broadcastRequestHandlers[m] = h
	mp.broadcastRequestHandlersMtx.Unlock()
}

func (mp *Meepo) getBroadcastRequestHandler(m Method) (BroadcastRequestHandler, error) {
	mp.broadcastRequestHandlersMtx.Lock()
	handler, ok := mp.broadcastRequestHandlers[m]
	mp.broadcastRequestHandlersMtx.Unlock()
	if !ok {
		return nil, ErrUnsupportedMethod
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

func (mp *Meepo) onSysDataChannelMessage(dc transport.DataChannel, b []byte) {
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onDataChannelMessageHandler",
		"label":   dc.Label(),
	})

	in, err := packet.UnmarshalPacket(b)
	if err != nil {
		logger.WithError(err).Debugf("failed to unmarshal buffer")
		return
	}

	switch in.Header().Type() {
	case packet.Request:
		mp.handleRequest(dc, in)
	case packet.Response:
		mp.handleResponse(dc, in)
	case packet.BroadcastRequest:
		mp.handleBroadcastRequest(dc, in)
	case packet.BroadcastResponse:
		mp.handleBroadcastResponse(dc, in)
	}
}

func (mp *Meepo) recoverHandleRequest(p packet.Packet) {
	recovered := recover()
	if recovered == nil {
		return
	}

	hdr := p.Header()
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":     "recoverHandleRequest",
		"method":      hdr.Method(),
		"session":     hdr.Session(),
		"source":      hdr.Source(),
		"destination": hdr.Destination(),
	})

	switch recovered.(type) {
	case errSendPacket:
		err := recovered.(errSendPacket)
		logger.WithError(err).Warningf("send packet error")
	case error:
		err := recovered.(error)
		logger.WithError(err).Debugf("failed to handle request")
	}
}

func (mp *Meepo) initRequestHandlers() {
	mp.registerRequestHandleFunc(METHOD_PING, mp.onPing)
	mp.registerRequestHandleFunc(METHOD_NEW_TELEPORTATION, mp.onNewTeleportation)
	mp.registerRequestHandleFunc(METHOD_CLOSE_TELEPORTATION, mp.onCloseTeleportation)
	mp.registerRequestHandleFunc(METHOD_DO_TELEORT, mp.onDoTeleport)
}

func (mp *Meepo) handleRequest(dc transport.DataChannel, in packet.Packet) {
	hdr := in.Header()

	defer mp.recoverHandleRequest(in)

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":     "handleRequest",
		"method":      hdr.Method(),
		"session":     hdr.Session(),
		"source":      hdr.Source(),
		"destination": hdr.Destination(),
	})

	if err := mp.authenticatePacket(in); err != nil {
		logger.WithError(err).Debugf("unauthenticated request")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	handler, err := mp.getRequestHandler(Method(hdr.Method()))
	if err != nil {
		logger.WithError(err).Debugf("failed to get request handler")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	handler.Handle(dc, in)
	logger.Tracef("done")
}

func (mp *Meepo) handleResponse(dc transport.DataChannel, in packet.Packet) {
	hdr := in.Header()
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":     "handleResponse",
		"source":      hdr.Source(),
		"destination": hdr.Destination(),
		"method":      hdr.Method(),
		"session":     hdr.Session(),
	})

	ch, unlock, err := mp.channelLocker.GetWithUnlock(hdr.Session())
	if err != nil {
		logger.WithError(err).Debugf("failed to get session channel")
		return
	}
	defer unlock()

	ch <- in
	logger.Tracef("send response to session channel")
}

func (mp *Meepo) sendRequest(in packet.Packet) error {
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":     "sendRequest",
		"destination": in.Header().Destination(),
		"session":     in.Header().Session(),
		"method":      in.Header().Method(),
	})

	tp, err := mp.getTransport(in.Header().Destination())
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

	mp.sendPacket(dc, in)

	return nil
}

func (mp *Meepo) sendPacket(dc transport.DataChannel, p packet.Packet) {
	buf, err := packet.MarshalPacket(p)
	if err != nil {
		panic(ErrSendPacket(err))
	}

	if _, err = dc.Write(buf); err != nil {
		panic(ErrSendPacket(err))
	}
}

func (mp *Meepo) waitResponse(session int32) (packet.Packet, error) {
	ch, err := mp.channelLocker.Get(session)
	if err != nil {
		return nil, err
	}

	select {
	case out, ok := <-ch:
		if !ok {
			return nil, SessionChannelClosedError(session)
		}

		res := out.(packet.Packet)

		if err = res.Err(); err != nil {
			return nil, err
		}

		return res, nil
	case <-time.After(cast.ToDuration(mp.opt.Get("waitResponseTimeout").Inter())):
		return nil, ErrWaitResponseTimeout
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

func (mp *Meepo) doRequest(in packet.Packet) (packet.Packet, error) {
	var err error

	hdr := in.Header()
	dst := hdr.Destination()
	sess := hdr.Session()
	meth := hdr.Method()

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":     "doRequest",
		"destination": dst,
		"session":     sess,
		"method":      meth,
	})

	if in, err = mp.signPacket(in); err != nil {
		logger.WithError(err).Debugf("failed to sign message")
		return nil, err
	}

	if err = mp.channelLocker.Acquire(in.Header().Session()); err != nil {
		logger.WithError(err).Debugf("failed to acquire session channel")
		return nil, err
	}
	defer func() {
		mp.channelLocker.Release(sess)
		logger.Tracef("release session channel")
	}()
	logger.Tracef("acquire session channel")

	if err = mp.sendRequest(in); err != nil {
		logger.WithError(err).Debugf("failed to send request")
		return nil, err
	}
	logger.Tracef("send request")

	out, err := mp.waitResponse(sess)
	if err != nil {
		logger.WithError(err).Debugf("failed to wait response")
		return nil, err
	}
	logger.Tracef("receive response")

	if err = mp.authenticatePacket(out, WithSubject(dst)); err != nil {
		logger.WithError(err).Debugf("unauthenticated response")
		return nil, err
	}

	logger.Tracef("done")

	return out, nil
}

func (mp *Meepo) sendResponse(dc transport.DataChannel, out packet.Packet) {
	var err error

	if out, err = mp.signPacket(out); err != nil {
		panic(err)
	}

	mp.sendPacket(dc, out)
}
