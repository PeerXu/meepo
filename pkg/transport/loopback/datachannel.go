package loopback_transport

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/PeerXu/meepo/pkg/transport"
	"github.com/sirupsen/logrus"
)

type LoopbackDataChannel struct {
	ch1 chan []byte
	ch2 chan []byte

	left  *LoopbackDataChannelWrapper
	right *LoopbackDataChannelWrapper

	label     string
	transport transport.Transport
	state     transport.DataChannelState
	mtx       sync.Locker

	logger logrus.FieldLogger
}

func NewLoopbackDataChannel(label string, tp transport.Transport, logger logrus.FieldLogger) *LoopbackDataChannel {
	return &LoopbackDataChannel{
		ch1:       make(chan []byte),
		ch2:       make(chan []byte),
		label:     label,
		transport: tp,
		logger:    logger,
		state:     transport.DataChannelStateConnecting,
		mtx:       new(sync.Mutex),
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
	time.Sleep(1000 * time.Millisecond)

	ldc.mtx.Lock()
	defer ldc.mtx.Unlock()

	if ldc.state == transport.DataChannelStateClosed {
		return nil
	}

	ldc.state = transport.DataChannelStateClosing
	logger.WithField("state", transport.DataChannelStateClosing).Tracef("change state")

	close(ldc.ch1)
	logger.WithField("channel", "ch1").Tracef("inner channel closed")
	close(ldc.ch2)
	logger.WithField("channel", "ch2").Tracef("inner channel closed")

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
			ldc, ldc.ch1, ldc.ch2, ldc.getRawLogger().WithField("side", "left"))
	}
	return ldc.left
}

func (ldc *LoopbackDataChannel) Right() *LoopbackDataChannelWrapper {
	ldc.mtx.Lock()
	defer ldc.mtx.Unlock()
	if ldc.right == nil {
		ldc.right = NewLoopbackDataChannelWrapper(
			ldc, ldc.ch2, ldc.ch1, ldc.getRawLogger().WithField("side", "right"))
	}
	return ldc.right
}

type LoopbackDataChannelWrapper struct {
	*LoopbackDataChannel
	logger            logrus.FieldLogger
	reader            chan []byte
	writer            chan []byte
	buffer            bytes.Buffer
	mtx               sync.Locker
	onOpenHandler     func()
	onOpenHandlerOnce sync.Once
}

func NewLoopbackDataChannelWrapper(
	ldc *LoopbackDataChannel,
	reader chan []byte, writer chan []byte,
	logger logrus.FieldLogger,
) *LoopbackDataChannelWrapper {
	return &LoopbackDataChannelWrapper{
		LoopbackDataChannel: ldc,

		reader: reader,
		writer: writer,
		mtx:    new(sync.Mutex),
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
	if w.buffer.Len() == 0 {
		select {
		case buf, ok := <-w.reader:
			if !ok {
				return 0, io.EOF
			}

			w.buffer.Write(buf)
		}
	}

	return w.buffer.Read(p)
}

func (w *LoopbackDataChannelWrapper) Write(p []byte) (int, error) {
	if w.State() != transport.DataChannelStateOpen {
		return 0, io.EOF
	}
	w.writer <- p
	return len(p), nil
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
