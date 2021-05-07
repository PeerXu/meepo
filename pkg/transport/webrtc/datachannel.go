package webrtc_transport

import (
	"sync"
	"time"

	"github.com/pion/datachannel"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/transport"
)

type WebrtcDataChannel struct {
	logger logrus.FieldLogger

	dc  *webrtc.DataChannel
	ddc datachannel.ReadWriteCloser

	detachOnce sync.Once
	opened     chan struct{}

	transport *WebrtcTransport
}

func (wdc *WebrtcDataChannel) getLogger() logrus.FieldLogger {
	return wdc.logger.WithFields(logrus.Fields{
		"#instance": "WebrtcDataChannel",
		"label":     wdc.Label(),
	})
}

func (wdc *WebrtcDataChannel) waitOpened() error {
	select {
	case <-wdc.opened:
		return nil
	case <-time.After(5 * time.Second):
		return WaitDataChannelOpenedTimeoutError
	}
}

func (wdc *WebrtcDataChannel) Transport() transport.Transport {
	return wdc.transport
}

func (wdc *WebrtcDataChannel) Label() string {
	return wdc.dc.Label()
}

func (wdc *WebrtcDataChannel) State() transport.DataChannelState {
	return transport.DataChannelState(wdc.dc.ReadyState())
}

func (wdc *WebrtcDataChannel) OnOpen(f func()) {
	wdc.dc.OnOpen(func() {
		wdc.detachOnce.Do(func() {
			wdc.ddc, _ = wdc.dc.Detach()
			close(wdc.opened)
		})
		f()
	})
}

func (wdc *WebrtcDataChannel) Read(p []byte) (int, error) {
	if err := wdc.waitOpened(); err != nil {
		return 0, err
	}

	return wdc.ddc.Read(p)
}

func (wdc *WebrtcDataChannel) Write(p []byte) (int, error) {
	if err := wdc.waitOpened(); err != nil {
		return 0, err
	}

	return wdc.ddc.Write(p)
}

func (wdc *WebrtcDataChannel) Close() error {
	return wdc.dc.Close()
}

func NewWebrtcDataChannel(logger logrus.FieldLogger, dc *webrtc.DataChannel, tp *WebrtcTransport) *WebrtcDataChannel {
	wdc := &WebrtcDataChannel{
		logger:    logger,
		dc:        dc,
		opened:    make(chan struct{}),
		transport: tp,
	}

	return wdc
}
