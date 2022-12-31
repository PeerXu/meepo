package transport_webrtc

import (
	"math"
	"sync"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

const (
	datachannelBufferSize = math.MaxUint16
)

var rlBufPool = sync.Pool{New: func() any {
	return make([]byte, datachannelBufferSize)
}}

func (t *WebrtcTransport) readLoop() {
	logger := t.GetLogger().WithField("#method", "readLoop")

	for {
		buffer := rlBufPool.Get().([]byte)
		n, err := t.rwc.Read(buffer)
		if err != nil {
			rlBufPool.Put(buffer) // nolint:staticcheck
			logger.WithError(err).Debugf("read ReadWriteCloser error")
			return
		}
		data := make([]byte, n)
		copy(data, buffer[:n])
		rlBufPool.Put(buffer) // nolint:staticcheck
		go t.onMessage(data)
		logger.WithFields(logging.Fields{
			"bytes": n,
		}).Tracef("read message")
	}

}
