package loopback_transport

import (
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/transport"
	msync "github.com/PeerXu/meepo/pkg/util/sync"
)

type LoopbackDataChannel struct {
	c1Rd *io.PipeReader
	c1Wr *io.PipeWriter
	c2Rd *io.PipeReader
	c2Wr *io.PipeWriter

	left  *LoopbackDataChannelWrapper
	right *LoopbackDataChannelWrapper

	label     string
	transport transport.Transport
	state     transport.DataChannelState
	mtx       msync.Locker

	logger logrus.FieldLogger
}

func NewLoopbackDataChannel(label string, tp transport.Transport, logger logrus.FieldLogger) *LoopbackDataChannel {
	c1Rd, c1Wr := io.Pipe()
	c2Rd, c2Wr := io.Pipe()

	return &LoopbackDataChannel{
		c1Rd:      c1Rd,
		c1Wr:      c1Wr,
		c2Rd:      c2Rd,
		c2Wr:      c2Wr,
		label:     label,
		transport: tp,
		logger:    logger,
		state:     transport.DataChannelStateConnecting,
		mtx:       msync.NewLock(),
	}
}

func (ldc *LoopbackDataChannel) getLogger() logrus.FieldLogger {
	return ldc.logger.WithFields(logrus.Fields{
		"#instance": "LoopbackDataChannel",
		"label":     ldc.Label(),
	})
}

func (ldc *LoopbackDataChannel) getRawLogger() logrus.FieldLogger {
	return ldc.logger
}

func (ldc *LoopbackDataChannel) Transport() transport.Transport {
	ldc.mtx.Lock()
	defer ldc.mtx.Unlock()
	return ldc.transport
}

func (ldc *LoopbackDataChannel) Label() string {
	return ldc.label
}

func (ldc *LoopbackDataChannel) setState(s transport.DataChannelState) {
	logger := ldc.getLogger().WithField("#method", "setState")

	ldc.mtx.Lock()
	ldc.state = s
	ldc.mtx.Unlock()

	logger.WithField("state", s).Tracef("set state")
}

func (ldc *LoopbackDataChannel) State() transport.DataChannelState {
	ldc.mtx.Lock()
	defer ldc.mtx.Unlock()
	return ldc.state
}

func (ldc *LoopbackDataChannel) OnOpen(f func()) {
	panic("unimplemented")
}

func (ldc *LoopbackDataChannel) Read([]byte) (int, error) {
	panic("unimplemented")
}

func (ldc *LoopbackDataChannel) Write([]byte) (int, error) {
	panic("unimplemented")
}

func (ldc *LoopbackDataChannel) Close() error {
	logger := ldc.getLogger().WithField("#method", "Close")

	// HACK: lazy closing wait for send response to other side
	time.Sleep(50 * time.Millisecond)

	ldc.mtx.Lock()
	defer ldc.mtx.Unlock()

	if ldc.state == transport.DataChannelStateClosed {
		return nil
	}

	ldc.state = transport.DataChannelStateClosing
	logger.WithField("state", transport.DataChannelStateClosing).Tracef("change state")

	logger.WithError(ldc.c1Rd.Close()).Tracef("c1rd closed")
	logger.WithError(ldc.c1Wr.Close()).Tracef("c1wr closed")
	logger.WithError(ldc.c2Rd.Close()).Tracef("c2rd closed")
	logger.WithError(ldc.c2Wr.Close()).Tracef("c2wr closed")

	ldc.state = transport.DataChannelStateClosed
	logger.WithField("state", transport.DataChannelStateClosed).Tracef("change state")

	defer logger.Tracef("closed")

	return nil
}

func (ldc *LoopbackDataChannel) Left() *LoopbackDataChannelWrapper {
	ldc.mtx.Lock()
	defer ldc.mtx.Unlock()
	if ldc.left == nil {
		ldc.left = NewLoopbackDataChannelWrapper(
			ldc, ldc.c1Rd, ldc.c2Wr, ldc.getRawLogger().WithField("side", "left"))
	}
	return ldc.left
}

func (ldc *LoopbackDataChannel) Right() *LoopbackDataChannelWrapper {
	ldc.mtx.Lock()
	defer ldc.mtx.Unlock()
	if ldc.right == nil {
		ldc.right = NewLoopbackDataChannelWrapper(
			ldc, ldc.c2Rd, ldc.c1Wr, ldc.getRawLogger().WithField("side", "right"))
	}
	return ldc.right
}

type LoopbackDataChannelWrapper struct {
	*LoopbackDataChannel
	logger            logrus.FieldLogger
	reader            io.Reader
	writer            io.Writer
	mtx               msync.Locker
	onOpenHandler     func()
	onOpenHandlerOnce sync.Once
}

func NewLoopbackDataChannelWrapper(
	ldc *LoopbackDataChannel,
	reader io.Reader, writer io.Writer,
	logger logrus.FieldLogger,
) *LoopbackDataChannelWrapper {
	return &LoopbackDataChannelWrapper{
		LoopbackDataChannel: ldc,

		reader: reader,
		writer: writer,
		mtx:    msync.NewLock(),
		logger: logger,
	}
}

func (w *LoopbackDataChannelWrapper) getLogger() logrus.FieldLogger {
	return w.logger.WithFields(logrus.Fields{
		"#instance": "LoopbackDataChannelWrapper",
		"label":     w.Label(),
	})
}

func (w *LoopbackDataChannelWrapper) Read(p []byte) (int, error) {
	return w.reader.Read(p)
}

func (w *LoopbackDataChannelWrapper) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *LoopbackDataChannelWrapper) onOpen() {
	logger := w.getLogger().WithField("#method", "onOpen")

	w.mtx.Lock()
	handler := w.onOpenHandler
	w.mtx.Unlock()

	if handler != nil {
		w.onOpenHandlerOnce.Do(func() {
			go func() {
				handler()
				logger.Tracef("data channel opened")
			}()
		})
	}
}

func (w *LoopbackDataChannelWrapper) OnOpen(f func()) {
	logger := w.getLogger().WithField("#method", "OnOpen")

	w.mtx.Lock()
	w.onOpenHandlerOnce = sync.Once{}
	w.onOpenHandler = f
	w.mtx.Unlock()

	if w.State() == transport.DataChannelStateOpen {
		w.onOpenHandlerOnce.Do(func() {
			go func() {
				f()
				logger.Tracef("data channel opened")
			}()
		})
	}
}
