package meepo_core

import "context"

func (mp *Meepo) Close(ctx context.Context) error {
	logger := mp.GetLogger().WithField("#method", "Close")

	ts, err := mp.ListTransports(ctx)
	if err != nil {
		logger.WithError(err).Debugf("failed to list transports")
		return err
	}

	for _, t := range ts {
		if err = t.Close(ctx); err != nil {
			logger.WithError(err).Debugf("failed to close transport")
			return err
		}
	}

	mp.closeOnce.Do(func() {
		close(mp.closed)
		close(mp.naviRequests)
	})

	logger.Tracef("meepo closed")

	return nil
}

func (mp *Meepo) isClosed() bool {
	select {
	case <-mp.closed:
		return true
	default:
		return false
	}
}
