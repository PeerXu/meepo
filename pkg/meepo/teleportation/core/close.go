package teleportation_core

import "context"

func (tp *teleportation) Close(ctx context.Context) error {
	var err error
	logger := tp.GetLogger().WithField("#method", "Close")

	if err = tp.onClose(tp); err != nil {
		logger.WithError(err).Debugf("on close failed")
	}

	if err = tp.listener.Close(); err != nil {
		logger.WithError(err).Debugf("failed to close listener")
	}

	logger.Tracef("teleportation closed")

	return nil
}
