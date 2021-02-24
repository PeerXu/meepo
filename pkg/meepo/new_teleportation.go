package meepo

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/teleportation"
	"github.com/PeerXu/meepo/pkg/transport"
)

type NewTeleportationRequest struct {
	Message

	Name          string `json:"name"`
	LocalNetwork  string `json:"localNetwork"`
	LocalAddress  string `json:"localAddress"`
	RemoteNetwork string `json:"remoteNetwork"`
	RemoteAddress string `json:"remoteAddress"`
}

type NewTeleportationResponse struct {
	Message
}

type DoTeleportRequest struct {
	Message

	Name  string `json:"name"`
	Label string `json:"label"`
}

type DoTeleportResponse struct {
	Message
}

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
		Message: Message{
			PeerID:  mp.GetID(),
			Type:    "request",
			Session: random.Int31(),
			Method:  "newTeleportation",
		},
		Name:          name,
		LocalNetwork:  local.Network(),
		LocalAddress:  local.String(),
		RemoteNetwork: remote.Network(),
		RemoteAddress: remote.String(),
	}

	out, err := mp.doRequest(id, req)
	if err != nil {
		logger.WithError(err).Errorf("failed to do request")
		return nil, err
	}

	res := out.(*NewTeleportationResponse)
	if res.Error != "" {
		err = fmt.Errorf(res.Error)
		logger.WithError(err).Errorf("failed to new teleportation by peer")
		return nil, err
	}

	tp, err := mp.GetTransport(id)
	if err != nil {
		logger.WithError(err).Errorf("failed to get transport")
		return nil, err
	}

	ts, err = teleportation.NewTeleportationSource(
		teleportation.WithLogger(mp.getRawLogger()),
		teleportation.WithName(name),
		teleportation.WithSource(local),
		teleportation.WithSink(remote),
		teleportation.WithTransport(tp),
		teleportation.WithDoTeleportFunc(func(label string) error {
			innerLogger := mp.getLogger().WithFields(logrus.Fields{
				"#method": "doTeleportFunc",
				"peerID":  id,
				"name":    name,
				"laddr":   local.String(),
				"raddr":   remote.String(),
			})

			req := &DoTeleportRequest{
				Message: Message{
					PeerID:  mp.GetID(),
					Type:    "request",
					Session: random.Int31(),
					Method:  "doTeleport",
				},
				Name:  name,
				Label: label,
			}

			out, err := mp.doRequest(id, req)
			if err != nil {
				innerLogger.WithError(err).Errorf("failed to do request")
				return err
			}

			res := out.(*DoTeleportResponse)
			if res.Error != "" {
				err = fmt.Errorf(res.Error)
				innerLogger.WithError(err).Errorf("failed to do teleport by peer")
				return err
			}

			innerLogger.Tracef("do teleport")

			return nil
		}),
		teleportation.WithOnCloseHandler(func() {
			mp.removeTeleportationSource(ts.Name())
			logger.Tracef("remove teleportation source")
		}),
		teleportation.WithOnErrorHandler(func(err error) {
			mp.removeTeleportationSource(ts.Name())
			logger.Tracef("remove teleportation source")
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
	logger.Infof("new teleportation source")

	mp.addTeleportationSource(ts.Name(), ts)
	logger.Tracef("add teleportation source")

	return ts, nil
}

func (mp *Meepo) onNewTeleportation(dc transport.DataChannel, in interface{}) {
	var ts *teleportation.TeleportationSink

	req := in.(*NewTeleportationRequest)

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onNewTeleportation",
		"peerID":  req.PeerID,
		"name":    req.Name,
		"laddr":   req.LocalAddress,
		"raddr":   req.RemoteAddress,
	})

	tp, err := mp.GetTransport(req.PeerID)
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		mp.sendMessage(dc, mp.invertMessageWithError(req, err))
		return
	}

	source, err := net.ResolveTCPAddr(req.LocalNetwork, req.LocalAddress)
	if err != nil {
		logger.WithError(err).Debugf("failed to resolve source addr")
		mp.sendMessage(dc, mp.invertMessageWithError(req, err))
		return
	}

	sink, err := net.ResolveTCPAddr(req.RemoteNetwork, req.RemoteAddress)
	if err != nil {
		logger.WithError(err).Debugf("failed to resolve sink addr")
		mp.sendMessage(dc, mp.invertMessageWithError(req, err))
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
		mp.sendMessage(dc, mp.invertMessageWithError(req, err))
		return
	}
	tp.OnTransportState(transport.TransportStateFailed, func(hid transport.HandleID) {
		ts.Close()
		ts.Transport().UnsetOnTransportState(transport.TransportStateFailed, hid)
	})
	logger.Tracef("new teleportation sink")

	mp.addTeleportationSink(ts.Name(), ts)
	logger.Tracef("add teleportation sink")

	res := &NewTeleportationResponse{
		Message: mp.invertMessage(req.Message),
	}
	mp.sendMessage(dc, res)

	logger.Tracef("done")
}

func (mp *Meepo) onDoTeleport(dc transport.DataChannel, in interface{}) {
	var err error

	req := in.(*DoTeleportRequest)

	logger := mp.getLogger().WithFields(logrus.Fields{
		"#method": "onDoTeleport",
		"peerID":  req.PeerID,
		"name":    req.Name,
	})

	ts, ok := mp.getTeleportationSink(req.Name)
	if !ok {
		err = TeleportationNotExistError
		logger.WithError(err).Errorf("failed to get teleportation sink")
		mp.sendMessage(dc, mp.invertMessageWithError(req, err))
		return
	}

	if err = ts.OnDoTeleport(req.Label); err != nil {
		logger.WithError(err).Errorf("failed to do teleport")
		mp.sendMessage(dc, mp.invertMessageWithError(req, err))
		return
	}
	logger.Tracef("do teleport")

	res := &DoTeleportResponse{
		Message: mp.invertMessage(req.Message),
	}
	mp.sendMessage(dc, res)

	logger.Debugf("done")
}

func init() {
	registerDecodeMessageHelper("request", "newTeleportation", func() interface{} { return &NewTeleportationRequest{} })
	registerDecodeMessageHelper("response", "newTeleportation", func() interface{} { return &NewTeleportationResponse{} })
	registerDecodeMessageHelper("request", "doTeleport", func() interface{} { return &DoTeleportRequest{} })
	registerDecodeMessageHelper("response", "doTeleport", func() interface{} { return &DoTeleportResponse{} })
}
