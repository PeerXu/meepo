package teleportation_core

import "context"

func (tp *teleportation) Close(ctx context.Context) error {
	var err error
	logger := tp.GetLogger().WithField("#method", "Close")

	if h := tp.beforeCloseTeleportationHook; h != nil {
		if err = h(tp); err != nil {
			logger.WithError(err).Debugf("before close teleportation hook failed")
			return err
		}
	}

	if err = tp.listener.Close(); err != nil {
		logger.WithError(err).Debugf("failed to close listener")
	}

	if h := tp.afterCloseTeleportationHook; h != nil {
		h(tp)
	}

	logger.Tracef("teleportation closed")

	return nil
}
