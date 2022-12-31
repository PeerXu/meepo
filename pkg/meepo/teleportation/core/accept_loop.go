package teleportation_core

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (tp *teleportation) acceptLoop() {
	logger := tp.GetLogger().WithFields(logging.Fields{
		"#method": "acceptLoop",
	})
	defer func() {
		tp.Close(tp.context()) // nolint:errcheck
		logger.Tracef("accept loop closed")
	}()

	for {
		conn, err := tp.listener.Accept()
		if err != nil {
			logger.WithError(err).Debugf("failed to accept")
			return
		}

		go tp.onAccept(tp, conn)
	}
}
