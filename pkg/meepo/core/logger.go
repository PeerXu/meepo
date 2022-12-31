package meepo_core

import "github.com/PeerXu/meepo/pkg/internal/logging"

func (mp *Meepo) GetLogger() logging.Logger {
	return mp.logger.WithFields(logging.Fields{
		"#instance": "Meepo",
		"addr":      mp.Addr(),
	})
}

func (mp *Meepo) GetRawLogger() logging.Logger {
	return mp.logger
}
