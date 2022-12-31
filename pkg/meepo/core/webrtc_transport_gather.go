package meepo_core

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/addr"
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

func (mp *Meepo) gatherFunc(target addr.Addr, gtkFn getTrackersFunc, opt gatherOption) transport_webrtc.GatherFunc {
	return func(offer webrtc.SessionDescription) (answer webrtc.SessionDescription, err error) {
		logger := mp.GetLogger().WithFields(logging.Fields{
			"#method": "gatherFunc",
			"target":  target.String(),
		})

		req, err := mp.newNewTransportRequest(target, offer, opt)
		if err != nil {
			logger.WithError(err).Debugf("failed to new NewTransport request")
			return
		}

		if gtkFn == nil {
			gtkFn = func(target Addr) ([]Tracker, bool, error) { return mp.getNearestTrackers(target, mp.dhtAlpha, nil) }
		}

		tks, found, err := gtkFn(target)
		if err != nil {
			logger.WithError(err).Debugf("failed to get nearest trackers")
			return
		}
		logger = logger.WithFields(logging.Fields{
			"found":           found,
			"trackers.length": len(tks),
		})

		if len(tks) == 0 {
			err = ErrNoAvailableTrackers
			logger.WithError(err).Debugf("no available trackers")
			return
		}

		if found {
			tks = tks[:1]
		}

		done := make(chan struct{})
		answers := make(chan webrtc.SessionDescription)
		errs := make(chan error)
		defer func() {
			close(done)
			close(answers)
			close(errs)
		}()

		for _, tk := range tks {
			go func(tk Tracker) {
				logger := logger.WithField("tracker", tk.Addr().String())
				_answer, _err := tk.NewTransport(req)
				select {
				case <-done:
					logger.Tracef("gather already done")
					return
				default:
				}

				if _err != nil {
					logger.WithError(_err).Tracef("failed to new transport by tracker")
					errs <- _err
					return
				}
				answers <- _answer
				logger.Tracef("new transport by tracker")
			}(tk)
		}

		for i := 0; i < len(tks); i++ {
			select {
			case answer = <-answers:
				logger.Tracef("gather done")
				return answer, nil
			case err = <-errs:
			}
		}

		logger.WithError(err).Debugf("failed to gather")
		return
	}
}
