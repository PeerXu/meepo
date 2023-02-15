package meepo_core

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	transport_webrtc "github.com/PeerXu/meepo/pkg/meepo/transport/webrtc"
)

type getTrackersFunc func(Addr) (tks []Tracker, found bool, err error)

type gatherOption struct {
	EnableMux    bool
	MuxLabel     string
	MuxVer       int
	MuxBuf       int
	MuxStreamBuf int
	MuxKeepalive int
	MuxNocomp    bool

	EnableKcp      bool
	KcpLabel       string
	KcpPreset      string
	KcpCrypt       string
	KcpKey         string
	KcpMtu         int
	KcpSndwnd      int
	KcpRcvwnd      int
	KcpDataShard   int
	KcpParityShard int
}

func (mp *Meepo) genGatherFunc(target addr.Addr) transport_webrtc.GatherFunc {
	return func(sess transport_webrtc.Session, offer webrtc.SessionDescription) (answer webrtc.SessionDescription, err error) {
		logger := mp.GetLogger().WithFields(logging.Fields{
			"#method": "gatherFunc",
			"target":  target.String(),
			"session": sess.String(),
		})
		req, err := mp.newAddPeerConnectionRequest(target, sess, offer)
		if err != nil {
			logger.WithError(err).Debugf("failed to new AddPeerConnection request")
			return
		}

		out, err := mp.forwardRequest(mp.context(), target, req,
			func(tk Tracker, in *crypto_core.Packet) (any, error) { return tk.AddPeerConnection(in) },
			func(target Addr) ([]Tracker, bool, error) {
				return mp.getCloserTrackers(target, mp.dhtAlpha, []Addr{target})
			},
			logger)
		if err != nil {
			return
		}

		answer = out.(webrtc.SessionDescription)

		return
	}
}

func (mp *Meepo) genGatherOnNewFunc(target addr.Addr, gtksFn getTrackersFunc, opt gatherOption) transport_webrtc.GatherFunc {
	return func(sess transport_webrtc.Session, offer webrtc.SessionDescription) (answer webrtc.SessionDescription, err error) {
		logger := mp.GetLogger().WithFields(logging.Fields{
			"#method": "gatherOnNewFunc",
			"target":  target.String(),
			"session": sess.String(),
		})
		req, err := mp.newNewTransportRequest(target, sess, offer, opt)
		if err != nil {
			logger.WithError(err).Debugf("failed to new NewTransport request")
			return
		}
		if gtksFn == nil {
			gtksFn = func(target Addr) ([]Tracker, bool, error) { return mp.getCloserTrackers(target, mp.dhtAlpha, nil) }
		}

		out, err := mp.forwardRequest(mp.context(), target, req,
			func(tk Tracker, in *crypto_core.Packet) (any, error) { return tk.NewTransport(in) },
			gtksFn,
			logger)
		if err != nil {
			return
		}

		answer = out.(webrtc.SessionDescription)

		return
	}
}
