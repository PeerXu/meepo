package teleportation_core

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (tp *teleportation) GetLogger() logging.Logger {
	return tp.logger.WithFields(logging.Fields{
		"#instance":     "teleportation",
		"id":            tp.ID(),
		"addr":          tp.Addr().String(),
		"sourceNetwork": tp.sourceAddr.Network(),
		"sourceAddress": tp.sourceAddr.String(),
		"sinkNetwork":   tp.sinkAddr.Network(),
		"sinkAddress":   tp.sinkAddr.String(),
	})
}
