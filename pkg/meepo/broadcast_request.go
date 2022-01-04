package meepo

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/transport"
	mgroup "github.com/PeerXu/meepo/pkg/util/group"
)

const (
	MAX_HOP_LIMITED = 16
)

func (mp *Meepo) initBroadcastRequestHandlers() {
	if mp.opt.Get("asSignaling").Bool() {
		mp.registerBroadcastRequestHandleFunc(METHOD_WIRE, mp.onWire, &registerBroadcastRequestHandleFuncOption{
			NewGroupFunc: mgroup.NewAnyGroupFunc,
		})
	}
}

type registerBroadcastRequestHandleFuncOption struct {
	NewGroupFunc func() mgroup.Group
}

func (mp *Meepo) registerBroadcastRequestHandleFunc(m Method, h func(transport.DataChannel, packet.BroadcastPacket), opts ...*registerBroadcastRequestHandleFuncOption) {
	var opt *registerBroadcastRequestHandleFuncOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = &registerBroadcastRequestHandleFuncOption{
			NewGroupFunc: mgroup.NewRaceGroupFunc,
		}
	}

	mp.registerBroadcastRequestHandler(m, &MeepoBroadcastRequestHandler{
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

func (mp *Meepo) doBroadcastRequest(in packet.Packet, opts ...*doBroadcastRequestOption) (packet.Packet, error) {
	var opt *doBroadcastRequestOption
	var tps []transport.Transport
	var out interface{}
	var err error

	hdr := in.Header()
	dst := hdr.Destination()

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":     "doBroadcastRequest",
		"destination": dst,
	})

	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultDoBroadcastRequestOption()
	}

	if tp, err := mp.getConnectedTransport(dst); err == nil {
		tps = append(tps, tp)
	} else {
		if !errors.Is(err, ErrTransportNotExist) {
			logger.WithError(err).Debugf("failed to get connected transport")
			return nil, err
		}

		tps, err = mp.listConnectedTransports()
		if err != nil {
			logger.WithError(err).Debugf("failed to list connected transports")
			return nil, err
		}

		if len(tps) == 0 {
			err = ErrOutOfEdge
			logger.WithError(err).Debugf("no connected transports")
			return nil, err
		}
	}

	if in, err = mp.signPacket(in); err != nil {
		logger.WithError(err).Debugf("failed to sign packet")
		return nil, err
	}

	gg := opt.NewGroupFunc()
	for _, tp := range tps {
		ds := tp.PeerID()
		innerLogger := logger.WithField("broadcastDestination", ds)
		gg.Go(func() (interface{}, error) {
			tin := mp.createBroadcastRequest(ds, in)
			tout, err := mp.doRequest(packet.BroadcastPacketToPacket(tin))
			if err != nil {
				innerLogger.WithError(err).Debugf("failed to do broadcast request")
				return nil, err
			}

			bout, err := packet.PacketToBroadcastPacketE(tout)
			if err != nil {
				innerLogger.WithError(err).Debugf("failed to cast to broadcast response")
				return nil, err
			}

			if err = mp.authenticatePacket(bout.Packet(), WithSubject(hdr.Destination())); err != nil {
				innerLogger.WithError(err).Debugf("unauthenticated broadcast response")
				return nil, err
			}

			innerLogger.Tracef("do broadcast request done")
			return bout, nil
		})
	}

	if out, err = gg.Wait(); err != nil {
		return nil, err
	}
	logger.Tracef("done")

	return out.(packet.BroadcastPacket).Packet(), nil
}

func (mp *Meepo) hashBroadcastPacket(p packet.BroadcastPacket) string {
	return fmt.Sprintf("%v.%v", p.Header().Method(), p.Header().Session())
}

func (mp *Meepo) recoverHandleBroadcastRequest(p packet.Packet) {
	recovered := recover()
	if recovered == nil {
		return
	}

	hdr := p.Header()
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":     "recoverHandleBroadcastRequest",
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
		logger.WithError(err).Debugf("failed to handle broadcast request")
	}
}

func (mp *Meepo) handleBroadcastRequest(dc transport.DataChannel, in packet.Packet) {
	var err error

	defer mp.recoverHandleBroadcastRequest(in)

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":               "handleBroadcastRequest",
		"method":                in.Header().Method(),
		"broadcast.session":     in.Header().Session(),
		"broadcast.source":      in.Header().Source(),
		"broadcast.destination": in.Header().Destination(),
	})

	if err = mp.authenticatePacket(in); err != nil {
		logger.WithError(err).Debugf("unauthenticated broadcast request, drop it")
		return
	}

	bin, err := packet.PacketToBroadcastPacketE(in)
	if err != nil {
		logger.WithError(err).Debugf("failed to cast to broadcast request")
		return
	}

	hdr := bin.Packet().Header()
	bhdr := bin.Header()

	logger = logger.WithFields(logrus.Fields{
		"session":     hdr.Session(),
		"source":      hdr.Source(),
		"destination": hdr.Destination(),
		"hop":         bhdr.Hop(),
	})

	if err = mp.authenticatePacket(bin.Packet()); err != nil {
		logger.WithError(err).Debugf("unauthenticated broadcast request, drop it")
		return
	}

	bh := mp.hashBroadcastPacket(bin)
	if _, ok := mp.broadcastCache.Get(bh); ok {
		logger.Debugf("handled broadcast request, drop it")
		return
	}
	mp.broadcastCache.Add(bh, nil)

	handler, err := mp.getBroadcastRequestHandler(Method(hdr.Method()))
	if err != nil {
		logger.WithError(err).Debugf("failed to get broadcast request handler")
		return
	}

	if mp.GetID() == hdr.Destination() {
		handler.Handle(dc, bin)
	} else {
		handler.HandleBroadcast(dc, bin)
	}
	logger.Tracef("done")
}

func (mp *Meepo) handleBroadcastResponse(dc transport.DataChannel, in packet.Packet) {
	hdr := in.Header()
	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method":               "handleBroadcastResponse",
		"broadcast.source":      hdr.Source(),
		"broadcast.destination": hdr.Destination(),
		"broadcast.session":     hdr.Session(),
		"method":                hdr.Method(),
	})

	ch, unlock, err := mp.channelLocker.GetWithUnlock(hdr.Session())
	if err != nil {
		logger.WithError(err).Debugf("failed to get session channel")
		return
	}
	defer unlock()

	ch <- in
	logger.Tracef("send broadcast response to session channel")
}

type newBroadcastHandleFuncOption struct {
	NewGroupFunc func() mgroup.Group
}

func (mp *Meepo) newBroadcastHandleFunc(opt *newBroadcastHandleFuncOption) func(transport.DataChannel, packet.BroadcastPacket) {
	if opt == nil {
		opt = &newBroadcastHandleFuncOption{
			NewGroupFunc: mgroup.NewRaceGroupFunc,
		}
	}

	return func(dc transport.DataChannel, bin packet.BroadcastPacket) {
		var tp transport.Transport
		var tps []transport.Transport
		var err error

		bhdr := bin.Header()
		logger := mp.getLogger().WithFields(logrus.Fields{
			"#method":               "handleBroadcast",
			"broadcast.session":     bhdr.Session(),
			"broadcast.source":      bhdr.Source(),
			"broadcast.destination": bhdr.Destination(),
			"method":                bhdr.Method(),
			"broadcast.hop":         bhdr.Hop(),
		})

		if bhdr.Hop() <= 0 {
			err = ErrOutOfEdge
			mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
			logger.WithError(err).Debugf("out of edge")
			return
		}

		if tp, err = mp.getConnectedTransport(bin.Header().Destination()); err == nil {
			tps = append(tps, tp)
		} else {
			if !errors.Is(err, ErrTransportNotExist) {
				mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
				logger.WithError(err).Debugf("failed to get connected transport")
				return
			}
		}

		if len(tps) == 0 {
			if tps, err = mp.listConnectedTransports(); err != nil {
				mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
				logger.WithError(err).Debugf("failed to list connected transports")
				return
			}

			// drop upstream transport
			var ttps []transport.Transport
			for _, tp := range tps {
				if tp.PeerID() != bhdr.Source() && tp.PeerID() != mp.GetID() {
					ttps = append(ttps, tp)
				}
			}
			tps = ttps

			if len(tps) == 0 {
				err = ErrOutOfEdge
				mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
				logger.WithError(err).Debugf("no connected transports")
				return
			}
		}

		gg := opt.NewGroupFunc()
		for _, tp := range tps {
			ds := tp.PeerID()
			innerLogger := logger.WithField("downstream", ds)
			gg.Go(func() (interface{}, error) {
				tin := mp.repackBroadcastRequest(ds, bin)
				tout, err := mp.doRequest(packet.BroadcastPacketToPacket(tin))
				if err != nil {
					innerLogger.WithError(err).Debugf("failed to do request")
					return nil, err
				}

				bout, err := packet.PacketToBroadcastPacketE(tout)
				if err != nil {
					innerLogger.WithError(err).Debugf("failed to cast to broadcast response")
					return nil, err
				}

				if err = mp.authenticatePacket(bout.Packet(), WithSubject(bin.Packet().Header().Destination())); err != nil {
					innerLogger.WithError(err).Debugf("unauthenticated broadcast response")
					return nil, err
				}

				innerLogger.Tracef("do request")
				return bout, nil
			}, func(out interface{}, err error) {
				if err != nil {
					mp.sendBroadcastResponse(dc, mp.createBroadcastResponseWithError(bin, err))
					logger.WithError(err).Debugf("failed to do broadcast request")
					return
				}

				mp.sendBroadcastResponse(dc, mp.repackBroadcastResponse(bin, out.(packet.BroadcastPacket)), SkipPacketSigning())
				innerLogger.Tracef("do response")
			})
		}

		gg.Wait()
		logger.Debugf("done")
	}

}

type sendBroadcastResponseOption = ofn.OFN

func SkipPacketSigning() ofn.OFN {
	return func(o ofn.Option) {
		o["skipPacketSigning"] = true
	}
}

func defaultSendBroadcastResponseOption() ofn.Option {
	return ofn.NewOption(map[string]interface{}{
		"skipPacketSigning": false,
	})
}

func (mp *Meepo) sendBroadcastResponse(dc transport.DataChannel, out packet.BroadcastPacket, opts ...sendBroadcastResponseOption) {
	var p packet.Packet
	var sp packet.Packet
	var err error

	o := defaultSendBroadcastResponseOption()
	for _, opt := range opts {
		opt(o)
	}

	if !cast.ToBool(o.Get("skipPacketSigning").Inter()) {
		if p, err = mp.signPacket(out.Packet()); err != nil {
			panic(err)
		}
		out = out.SetPacket(p)
	}

	if sp, err = mp.signPacket(packet.BroadcastPacketToPacket(out)); err != nil {
		panic(err)
	}

	mp.sendPacket(dc, sp)
}
