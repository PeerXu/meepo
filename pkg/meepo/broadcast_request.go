package meepo

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/transport"
	mgroup "github.com/PeerXu/meepo/pkg/util/group"
)

const (
	MAX_HOP_LIMITED = 16
)

type BroadcastRequest interface {
	MessageGetter
	BroadcastGetter
	Copier
}

func (mp *Meepo) initBroadcastRequestHandlers() {
	if mp.opt.Get("asSignaling").Bool() {
		mp.registerBroadcastRequestHandleFunc("wire", mp.onWire, &registerBroadcastRequestHandleFuncOption{
			NewGroupFunc: mgroup.NewAnyGroupFunc,
		})
	}
}

type registerBroadcastRequestHandleFuncOption struct {
	NewGroupFunc func() mgroup.Group
}

func (mp *Meepo) registerBroadcastRequestHandleFunc(name string, h func(transport.DataChannel, interface{}), opts ...*registerBroadcastRequestHandleFuncOption) {
	var opt *registerBroadcastRequestHandleFuncOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = &registerBroadcastRequestHandleFuncOption{
			NewGroupFunc: mgroup.NewRaceGroupFunc,
		}
	}

	mp.registerBroadcastRequestHandler(name, &MeepoBroadcastRequestHandler{
		handle: h,
		handleBroadcast: mp.newBroadcastHandleFunc(&newBroadcastHandleFuncOption{
			NewGroupFunc: opt.NewGroupFunc,
		}),
	})
}

type doBroadcastRequestOption struct {
	NewGroupFunc mgroup.NewGroupFunc
}

func defaultDoBroadcastRequestOption() *doBroadcastRequestOption {
	return &doBroadcastRequestOption{
		NewGroupFunc: mgroup.NewRaceGroupFunc,
	}
}

func (mp *Meepo) doBroadcastRequest(destinationID string, in interface{}, opts ...*doBroadcastRequestOption) (interface{}, error) {
	var opt *doBroadcastRequestOption
	var tps []transport.Transport
	var dss []string
	var out interface{}
	var err error

	req := in.(BroadcastRequest)

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":       "doBroadcastRequest",
		"destinationID": destinationID,
	})

	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultDoBroadcastRequestOption()
	}

	if tp, err := mp.getConnectedTransport(destinationID); err == nil {
		tps = append(tps, tp)
	} else {
		if !errors.Is(err, TransportNotExistError) {
			logger.WithError(err).Debugf("failed to get connected transport")
			return nil, err
		}

		tps, err = mp.listConnectedTransports()
		if err != nil {
			logger.WithError(err).Debugf("failed to list connected transports")
			return nil, err
		}

		if len(tps) == 0 {
			err = ReachTransportEdgeError
			logger.WithError(err).Debugf("no connected transports")
			return nil, err
		}
	}

	for _, tp := range tps {
		dss = append(dss, tp.PeerID())
	}

	gg := opt.NewGroupFunc()
	for _, ds := range dss {
		dsc := ds
		gg.Go(func() (interface{}, error) {
			breq := req.Copy().(BroadcastRequest)
			breq.GetMessage().Session = generateSession()
			out, err := mp.doRequest(dsc, breq)
			if err != nil {
				logger.WithError(err).Debugf("failed to do request")
				return nil, err
			}
			logger.Tracef("do request done")
			return out, nil
		})
	}

	if out, err = gg.Wait(); err != nil {
		// logged in mgroup.Group.Go
		return nil, err
	}
	logger.Tracef("done")

	return out, nil
}

func (mp *Meepo) handleBroadcastRequest(dc transport.DataChannel, in interface{}) {
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "handleBroadcastRequest",
	})

	req, ok := in.(BroadcastRequest)
	if !ok {
		logger.Debugf("in not a BroadcastRequest")
		return
	}

	m := req.GetMessage()
	b := req.GetBroadcast()

	logger = logger.WithFields(logrus.Fields{
		"method":           m.Method,
		"session":          m.Session,
		"peerID":           m.PeerID,
		"broadcastSession": b.BroadcastSession,
		"destinationID":    b.DestinationID,
		"hop":              b.Hop,
	})

	bid := b.Identifier()

	if mp.opt.Get("asSignaling").Bool() {
		if _, ok := mp.broadcastCache.Get(bid); ok {
			logger.Debugf("handled broadcast request, drop it")
			return
		}
		mp.broadcastCache.Add(bid, nil)
	}

	handler, err := mp.getBroadcastRequestHandler(m.Method)
	if err != nil {
		logger.WithError(err).Debugf("failed to get broadcast request handler")
		return
	}

	if mp.GetID() == b.DestinationID {
		handler.Handle(dc, in)
	} else {
		handler.HandleBroadcast(dc, in)
	}
	logger.Tracef("done")
}

type newBroadcastHandleFuncOption struct {
	NewGroupFunc func() mgroup.Group
}

func (mp *Meepo) newBroadcastHandleFunc(opt *newBroadcastHandleFuncOption) func(transport.DataChannel, interface{}) {
	if opt == nil {
		opt = &newBroadcastHandleFuncOption{
			NewGroupFunc: mgroup.NewRaceGroupFunc,
		}
	}

	return func(dc transport.DataChannel, in interface{}) {
		var tp transport.Transport
		var tps []transport.Transport
		var dss []string
		var err error

		logger := mp.getLogger().WithFields(logrus.Fields{
			"#method": "handleBroadcast",
		})

		req := in.(BroadcastRequest)
		m := req.GetMessage()
		b := req.GetBroadcast()

		logger = logger.WithFields(logrus.Fields{
			"method":           m.Method,
			"session":          m.Session,
			"broadcastSession": b.BroadcastSession,
			"destinationID":    b.DestinationID,
			"hop":              b.Hop,
		})

		if b.Hop <= 0 {
			err = HopIsZeroError
			mp.sendMessage(dc, mp.createBroadcastResponseWithError(err, req))
			logger.WithError(err).Debugf("hop is zero")
			return
		}

		if b.DetectNextHop {
			if tp, err = mp.getConnectedTransport(b.DestinationID); err == nil {
				tps = append(tps, tp)
			} else {
				if !errors.Is(err, TransportNotExistError) {
					mp.sendMessage(dc, mp.createBroadcastResponseWithError(err, req))
					logger.WithError(err).Debugf("failed to get connected transport")
					return
				}
			}
		}

		if len(tps) == 0 {
			if tps, err = mp.listConnectedTransports(); err != nil {
				mp.sendMessage(dc, mp.createBroadcastResponseWithError(err, req))
				logger.WithError(err).Debugf("failed to list connected transports")
				return
			}

			// drop upstream transport
			var ttps []transport.Transport
			for _, tp := range tps {
				if tp.PeerID() != m.PeerID {
					ttps = append(ttps, tp)
				}
			}
			tps = ttps

			if len(tps) == 0 {
				err = ReachTransportEdgeError
				mp.sendMessage(dc, mp.createBroadcastResponseWithError(err, req))
				logger.WithError(err).Debugf("no connected transports")
				return
			}
		}

		for _, tp = range tps {
			dss = append(dss, tp.PeerID())
		}

		gg := opt.NewGroupFunc()
		for _, ds := range dss {
			dsi := ds
			innerLogger := logger.WithField("downstream", dsi)
			gg.Go(func() (interface{}, error) {
				out, err := mp.doRequest(dsi, mp.createNextHopBroadcastRequest(req))
				if err != nil {
					innerLogger.WithError(err).Debugf("failed to do request")
					return nil, err
				}

				innerLogger.Tracef("do request")
				return out, nil
			}, func(out interface{}, err error) {
				if err != nil {
					mp.sendMessage(dc, mp.createBroadcastResponseWithError(err, req))
					logger.WithError(err).Debugf("failed to do broadcast request")
					return
				}

				mp.sendMessage(dc, mp.createBroadcastResponse(out, req))
				innerLogger.Tracef("do response")
			})
		}

		gg.Wait()
		logger.Debugf("done")
	}

}
