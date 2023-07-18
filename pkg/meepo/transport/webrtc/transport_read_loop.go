package transport_webrtc

import (
	"io"
	"math"
	"sync"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

const (
	datachannelBufferSize = math.MaxUint16
)

var rlBufPool = sync.Pool{New: func() any {
	return make([]byte, datachannelBufferSize)
}}

func (t *WebrtcTransport) readLoop(sess Session, rwc io.ReadWriteCloser) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "readLoop",
		"session": sess.String(),
	})

	for {
		buffer := rlBufPool.Get().([]byte)
		n, err := rwc.Read(buffer)
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