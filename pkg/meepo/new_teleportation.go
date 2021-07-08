package meepo

import (
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/meepo/auth"
	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/PeerXu/meepo/pkg/teleportation"
	"github.com/PeerXu/meepo/pkg/transport"
)

const (
	METHOD_NEW_TELEPORTATION Method = "newTeleportation"
	METHOD_DO_TELEORT        Method = "doTeleport"
)

type (
	NewTeleportationRequest struct {
		Name          string
		LocalNetwork  string
		LocalAddress  string
		RemoteNetwork string
		RemoteAddress string
		HashedSecret  string
	}

	NewTeleportationResponse struct{}

	DoTeleportRequest struct {
		Name  string
		Label string
	}

	DoTeleportResponse struct{}
)

func newNewTeleportationOption() objx.Map {
	return objx.New(map[string]interface{}{})
}

func (mp *Meepo) NewTeleportation(id string, remote net.Addr, opts ...NewTeleportationOption) (teleportation.Teleportation, error) {
	var ts *teleportation.TeleportationSource
	var local net.Addr
	var name string
	var ok bool
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "NewTeleportation",
		"peerID":  id,
	})

	o := newNewTeleportationOption()
	for _, opt := range opts {
		opt(o)
	}

	if local, ok = o.Get("local").Inter().(net.Addr); ok {
		if local, err = checkAddrIsListenable(local); err != nil {
			logger.WithError(err).Errorf("failed to check address")
			return nil, err
		}
	} else {
		local = getListenableAddr()
	}
	logger = logger.WithFields(logrus.Fields{
		"laddr": local.String(),
		"raddr": remote.String(),
	})

	if val := o.Get("name").Inter(); val == nil {
		o.Set("name", fmt.Sprintf("%s:%s", remote.Network(), remote.String()))
	}
	name = cast.ToString(o.Get("name").Inter())

	logger = logger.WithField("name", name)

	req := &NewTeleportationRequest{
		Name:          name,
		LocalNetwork:  local.Network(),
		LocalAddress:  local.String(),
		RemoteNetwork: remote.Network(),
		RemoteAddress: remote.String(),
	}

	secret := cast.ToString(o.Get("secret").Inter())
	if secret != "" {
		req.HashedSecret, err = hashSecret(secret)
		if err != nil {
			logger.WithError(err).Errorf("failed to hash secret")
			return nil, err
		}
	}

	in := mp.createRequest(id, METHOD_NEW_TELEPORTATION, req)

	out, err := mp.doRequest(in)
	if err != nil {
		logger.WithError(err).Errorf("failed to do request")
		return nil, err
	}

	if err = out.Err(); err != nil {
		logger.WithError(err).Errorf("failed to new teleportation by peer")
		return nil, err
	}

	tp, err := mp.getTransport(id)
	if err != nil {
		logger.WithError(err).Errorf("failed to get transport")
		return nil, err
	}

	var lisCloseOnce sync.Once
	dialRequests := make(chan *teleportation.DialRequest)
	lis, err := net.Listen(local.Network(), local.String())
	if err != nil {
		logger.WithError(err).Errorf("failed to listen local address")
		return nil, err
	}
	lisCloser := func() {
		logger.WithError(lis.Close()).Tracef("listener closed")
	}

	go func() {
		innerLogger := mp.getLogger().WithFields(logrus.Fields{
			"#method": "accpetLoop",
		})

		defer close(dialRequests)
		defer lisCloseOnce.Do(lisCloser)
		for {
			conn, err := lis.Accept()
			if err != nil {
				innerLogger.WithError(err).Debugf("failed to accept from listener")
				return
			}
			dialRequests <- teleportation.NewDialRequest(conn)
			innerLogger.Tracef("accepted")
		}
	}()

	ts, err = teleportation.NewTeleportationSource(
		teleportation.WithLogger(mp.getRawLogger()),
		teleportation.WithName(name),
		teleportation.WithSource(local),
		teleportation.WithSink(remote),
		teleportation.WithTransport(tp),
		teleportation.SetDialRequestChannel(dialRequests),
		teleportation.WithDoTeleportFunc(func(label string) error {
			innerLogger := mp.getLogger().WithFields(logrus.Fields{
				"#method": "doTeleportFunc",
				"peerID":  id,
				"name":    name,
				"laddr":   local.String(),
				"raddr":   remote.String(),
			})

			req := &DoTeleportRequest{
				Name:  name,
				Label: label,
			}

			in := mp.createRequest(id, METHOD_DO_TELEORT, req)

			out, err := mp.doRequest(in)
			if err != nil {
				innerLogger.WithError(err).Errorf("failed to do request")
				return err
			}

			if err = out.Err(); err != nil {
				innerLogger.WithError(err).Errorf("failed to do teleport by peer")
				return err
			}

			innerLogger.Tracef("do teleport")

			return nil
		}),
		teleportation.WithOnCloseHandler(func() {
			mp.removeTeleportationSource(ts.Name())
			logger.Tracef("remove teleportation source")

			lisCloseOnce.Do(lisCloser)
		}),
		teleportation.WithOnErrorHandler(func(err error) {
			mp.removeTeleportationSource(ts.Name())
			logger.WithError(err).Tracef("remove teleportation source")

			lisCloseOnce.Do(lisCloser)
		}),
	)
	if err != nil {
		logger.WithError(err).Errorf("failed to new teleportation source")
		return nil, err
	}

	tp.OnTransportState(transport.TransportStateFailed, func(hid transport.HandleID) {
		ts.Close()
		tp.UnsetOnTransportState(transport.TransportStateFailed, hid)
	})

	mp.addTeleportationSource(ts.Name(), ts)
	logger.Tracef("add teleportation source")

	logger.Infof("new teleportation source")

	return ts, nil
}

func (mp *Meepo) onNewTeleportation(dc transport.DataChannel, in packet.Packet) {
	var ts *teleportation.TeleportationSink
	var req NewTeleportationRequest
	var err error

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onNewTeleportation",
		"peerID":  in.Header().Source(),
	})

	if err = in.Data(&req); err != nil {
		logger.WithError(err).Errorf("failed to unmarshal request data")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"name":  req.Name,
		"laddr": req.LocalAddress,
		"raddr": req.RemoteAddress,
	})

	var authOpts []auth.AuthorizeOption
	if req.HashedSecret == "" {
		authOpts = append(
			authOpts,
			WithAuthorizationName("dummy"),
		)
	} else {
		authOpts = append(
			authOpts,
			WithAuthorizationName("secret"),
			WithAuthorizationSecret(req.HashedSecret),
		)
	}

	if err = mp.authorization.Authorize(in.Header().Source(), mp.GetID(), string(METHOD_NEW_TELEPORTATION), authOpts...); err != nil {
		logger.WithError(err).Debugf("unauthorized request")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	tp, err := mp.getTransport(in.Header().Source())
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	source, err := mp.resolveTeleportationSourceAddr(req.LocalNetwork, req.LocalAddress)
	if err != nil {
		logger.WithError(err).Debugf("failed to resolve source addr")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	sink, err := mp.resolveTeleportationSinkAddr(req.RemoteNetwork, req.RemoteAddress)
	if err != nil {
		logger.WithError(err).Debugf("failed to resolve sink addr")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	challenge := NewAclChallenge(NewAclEntity(in.Header().Source(), source), NewAclEntity(mp.GetID(), sink))
	logger = logger.WithField("challenge", challenge)
	if err = mp.acl.Allowed(challenge); err != nil {
		logger.WithError(err).Debugf("not allowed by acl")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	ts, err = teleportation.NewTeleportationSink(
		teleportation.WithLogger(mp.getRawLogger()),
		teleportation.WithName(req.Name),
		teleportation.WithSource(source),
		teleportation.WithSink(sink),
		teleportation.WithTransport(tp),
		teleportation.WithOnCloseHandler(func() {
			mp.removeTeleportationSink(ts.Name())
			logger.Tracef("remove teleportation sink")
		}),
		teleportation.WithOnErrorHandler(func(err error) {
			mp.removeTeleportationSink(ts.Name())
			logger.Tracef("remove teleportation sink")
		}),
	)
	if err != nil {
		logger.WithError(err).Debugf("failed to new teleportation sink")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}
	tp.OnTransportState(transport.TransportStateFailed, func(hid transport.HandleID) {
		ts.Close()
		ts.Transport().UnsetOnTransportState(transport.TransportStateFailed, hid)
	})
	logger.Tracef("new teleportation sink")

	mp.addTeleportationSink(ts.Name(), ts)
	logger.Tracef("add teleportation sink")

	mp.sendResponse(dc, mp.createResponse(in, &NewTeleportationResponse{}))
	logger.Tracef("done")
}

func (mp *Meepo) onDoTeleport(dc transport.DataChannel, in packet.Packet) {
	var err error
	var req DoTeleportRequest

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onDoTeleport",
		"peerID":  in.Header().Source(),
		"name":    req.Name,
	})

	if err = in.Data(&req); err != nil {
		logger.WithError(err).Errorf("failed to unmarshal request data")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	ts, ok := mp.getTeleportationSink(req.Name)
	if !ok {
		err = ErrTeleportationNotExist
		logger.WithError(err).Errorf("failed to get teleportation sink")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}

	if err = ts.OnDoTeleport(req.Label); err != nil {
		logger.WithError(err).Errorf("failed to do teleport")
		mp.sendResponse(dc, mp.createResponseWithError(in, err))
		return
	}
	logger.Tracef("do teleport")

	mp.sendResponse(dc, mp.createResponse(in, &DoTeleportResponse{}))
	logger.Debugf("done")
}
