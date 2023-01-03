package lock

import "github.com/PeerXu/meepo/pkg/lib/logging"

func (x *tracedLock) GetLogger() logging.Logger {
	return x.logger.WithFields(logging.Fields{
		"#instance":   "lock",
		"id":          x.id,
		"name":        x.name,
		"goroutineID": goid(),
	})
}
