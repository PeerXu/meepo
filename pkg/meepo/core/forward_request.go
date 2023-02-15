package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) forwardRequest(
	ctx context.Context,
	target addr.Addr,
	req *crypto_core.Packet,
	doReqFn func(Tracker, *crypto_core.Packet) (any, error),
	gtksFn getTrackersFunc,
	logger logging.Logger,
) (any, error) {
	tks, found, err := gtksFn(target)
	if err != nil {
		logger.WithError(err).Debugf("failed to get trackers")
		return nil, err
	}
	logger = logger.WithField("found", found)

	if len(tks) == 0 {
		err = ErrNoAvailableTrackers
		logger.WithError(err).Debugf("failed to get trackers")
	}

	if found {
		tks = tks[:1]
	}

	done := make(chan struct{}, 1)
	ress := make(chan any)
	errs := make(chan error)
	defer func() {
		close(done)
		close(ress)
		close(errs)
	}()

	for _, tk := range tks {
		go func(tk Tracker) {
			logger := logger.WithField("tracker", tk.Addr().String())
			res, err := doReqFn(tk, req)
			select {
			case <-done:
				logger.Tracef("forward already done")
				return
			default:
			}

			if err != nil {
				logger.WithError(err).Tracef("failed to forward request")
				errs <- err
				return
			}

			ress <- res
			logger.Tracef("forward request success")
		}(tk)
	}

	for i := 0; i < len(tks); i++ {
		select {
		case res := <-ress:
			logger.Tracef("forward done")
			return res, nil
		case err = <-errs:
		}
	}

	logger.WithError(err).Debugf("failed to forward request")
	return nil, err
}
