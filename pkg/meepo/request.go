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

func (mp *Meepo) registerSessionChannel(session int32) error {
	ch := make(chan interface{})

	if _, loaded := mp.sessionChannels.LoadOrStore(session, ch); loaded {
		defer close(ch)
		return SessionChannelExistError(session)
	}

	return nil
}

func (mp *Meepo) unregisterSessionChannel(session int32) error {
	ch, err := mp.getSessionChannel(session)
	if err != nil {
		return err
	}
	defer close(ch)

	mp.sessionChannels.Delete(session)

	return nil
}

func (mp *Meepo) getSessionChannel(session int32) (chan interface{}, error) {
	ch, ok := mp.sessionChannels.Load(session)
	if !ok {
		return nil, SessionChannelNotExistError(session)
	}

	return ch.(chan interface{}), nil
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
	case "request":
		mp.handleRequest(dc, in)
	case "response":
		if err = mp.handleResponse(dc, in); err != nil {
			logger.WithError(err).Debugf("failed to handle response")
		}
	}
}

func (mp *Meepo) recoverHandleRequest(m Message) {
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

func (mp *Meepo) handleRequest(dc transport.DataChannel, in interface{}) error {
	var handle func(transport.DataChannel, interface{})
	var err error

	m := in.(MessageGetter).GetMessage()

	defer mp.recoverHandleRequest(m)

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "handleRequest",
		"method":  m.Method,
		"session": m.Session,
	})

	switch m.Method {
	case "ping":
		handle = mp.onPing
	case "newTeleportation":
		handle = mp.onNewTeleportation
	case "closeTeleportation":
		handle = mp.onCloseTeleportation
	case "doTeleport":
		handle = mp.onDoTeleport
	}

	if handle == nil {
		err = UnsupportedRequestHandlerError(m.Method)
		logger.WithError(err).Debugf("unsupported request handler")
		return err
	}
	handle(dc, in)

	logger.Tracef("done")

	return nil
}

func (mp *Meepo) handleResponse(dc transport.DataChannel, in interface{}) error {
	m := in.(MessageGetter).GetMessage()

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "handleResponse",
		"method":  m.Method,
		"session": m.Session,
	})

	ch, err := mp.getSessionChannel(m.Session)
	if err != nil {
		logger.WithError(err).Debugf("failed to get session channel")
		return err
	}

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

	tp, err := mp.GetTransport(id)
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
	ch, err := mp.getSessionChannel(session)
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

	if err = mp.registerSessionChannel(msg.Session); err != nil {
		logger.WithError(err).Debugf("failed to register session channel")
		return nil, err
	}
	defer func() {
		mp.unregisterSessionChannel(msg.Session)
		logger.Tracef("unregister session channel")
	}()
	logger.Tracef("register session channel")

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
