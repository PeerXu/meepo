package transport_webrtc

import (
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/xtaci/smux"

	mio "github.com/PeerXu/meepo/pkg/internal/io"
	"github.com/PeerXu/meepo/pkg/internal/logging"
)

type tempDataChannel struct {
	dc  *webrtc.DataChannel
	rwc io.ReadWriteCloser
	req *NewChannelRequest
}

func (t *WebrtcTransport) onDataChannel(dc *webrtc.DataChannel) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "onDataChannel",
		"label":   dc.Label(),
	})

	if t.isSysDataChannel(dc) {
		dc.OnOpen(func() { t.onSysDataChannelOpen(dc) })
		logger.Tracef("setup sys data channel")
		return
	}

	if t.isMuxDataChannel(dc) {
		dc.OnOpen(func() { t.onMuxDataChannelOpen(dc) })
		logger.Tracef("setup mux data channel")
		return
	}

	if t.isKcpDataChannel(dc) {
		dc.OnOpen(func() { t.onKcpDataChannelOpen(dc) })
		logger.Tracef("setup kcp data channel")
		return
	}

	t.tempDataChannelsMtx.Lock()
	defer t.tempDataChannelsMtx.Unlock()

	dc.OnOpen(t.onDataChannelOpen(dc))

	tdc, found := t.tempDataChannels[dc.Label()]
	if !found {
		t.tempDataChannels[dc.Label()] = &tempDataChannel{dc: dc}
		go t.removeTimeoutTempDataChannel(dc.Label())
		logger.Tracef("create temp data channel")
		return
	}

	tdc.dc = dc
	logger.Tracef("assign webrtc.DataChannel to temp data channel")
}

func (t *WebrtcTransport) onDataChannelOpen(dc *webrtc.DataChannel) func() {
	return func() {
		var err error

		logger := t.GetLogger().WithFields(logging.Fields{
			"#method": "onDataChannelOpen",
			"label":   dc.Label(),
		})

		t.tempDataChannelsMtx.Lock()
		defer t.tempDataChannelsMtx.Unlock()

		tdc, found := t.tempDataChannels[dc.Label()]
		if !found {
			logger.Debugf("temp data channel not found")
			return
		}

		if tdc.rwc, err = dc.Detach(); err != nil {
			logger.WithError(err).Debugf("failed to detach data channel")
			return
		}

		if tdc.req == nil {
			logger.Tracef("wait for new channel request")
			return
		}

		go t.handleNewChannel(dc.Label())
		logger.Tracef("on data channel open")
	}
}

func (t *WebrtcTransport) removeTimeoutTempDataChannel(label string) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "removeTimeoutTempDataChannel",
		"label":   label,
	})

	<-time.After(t.tempDataChannelTimeout)

	t.tempDataChannelsMtx.Lock()
	defer t.tempDataChannelsMtx.Unlock()

	tdc, found := t.tempDataChannels[label]
	if !found {
		return
	}
	delete(t.tempDataChannels, label)

	if tdc.rwc != nil {
		if err := tdc.rwc.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close temp rwc")
		}
		tdc.rwc = nil
	}

	if tdc.dc != nil {
		if err := tdc.dc.Close(); err != nil {
			logger.WithError(err).Debugf("failed to close temp data channel")
		}
		tdc.dc = nil
	}

	logger.Tracef("remove timeout data channel")
}

func (t *WebrtcTransport) isSysDataChannel(dc *webrtc.DataChannel) bool {
	return strings.HasPrefix(dc.Label(), "sys#")
}

func (t *WebrtcTransport) isMuxDataChannel(dc *webrtc.DataChannel) bool {
	if dc.Label() == t.muxLabel {
		return true
	} else {
		return false
	}
}

func (t *WebrtcTransport) isKcpDataChannel(dc *webrtc.DataChannel) bool {
	if dc.Label() == t.kcpLabel {
		return true
	} else {
		return false
	}
}

func (t *WebrtcTransport) onSysDataChannelOpen(dc *webrtc.DataChannel) {
	var err error

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "onSysDataChannelOpen",
		"label":   dc.Label(),
	})

	t.sysMtx.Lock()
	defer t.sysMtx.Unlock()

	if t.sysDataChannel != nil && t.sysDataChannel.Label() != dc.Label() {
		logger.WithField("old", t.sysDataChannel.Label()).Debugf("sys data channel existed, close new one")
		dc.Close()
		return
	}

	rwc, err := dc.Detach()
	if err != nil {
		logger.WithError(err).Debugf("failed to detach DataChannel")
		dc.Close()
		return
	}

	t.rwc = rwc
	t.sysDataChannel = dc

	go t.readLoop()
	t.channelDone(1)
	logger.Tracef("on sys data channel open")
}

func (t *WebrtcTransport) onMuxDataChannelOpen(dc *webrtc.DataChannel) {
	var err error

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":       "onMuxDataChannelOpen",
		"label":         dc.Label(),
		"smuxVer":       t.muxVer,
		"smuxBuf":       t.muxBuf,
		"smuxStreamBuf": t.muxStreamBuf,
		"smuxKeepalive": t.muxKeepalive,
		"smuxNocomp":    t.muxNocomp,
	})

	t.muxMtx.Lock()
	defer t.muxMtx.Unlock()

	rwc, err := dc.Detach()
	if err != nil {
		logger.WithError(err).Debugf("failed to detach DataChannel")
		dc.Close()
		return
	}

	// HACK: compatibility issue: smux and pion webrtc
	pc1, pc2 := net.Pipe()

	done1 := make(chan struct{})
	go func() {
		defer close(done1)
		n, err := mio.Copy(pc2, rwc)
		logger.WithError(err).WithFields(logging.Fields{
			"from":  "webrtc.DataChannel",
			"to":    "net.Pipe",
			"bytes": n,
		}).Debugf("copy closed")
	}()

	done2 := make(chan struct{})
	go func() {
		defer close(done2)
		n, err := mio.Copy(rwc, pc2)
		logger.WithError(err).WithFields(logging.Fields{
			"from":  "net.Pipe",
			"to":    "webrtc.DataChannel",
			"bytes": n,
		}).Debugf("copy closed")
	}()

	go func() {
		select {
		case <-done1:
		case <-done2:
		}
		pc1.Close() // nolint:errcheck
		pc2.Close() // nolint:errcheck
	}()

	var newSessFn func(io.ReadWriteCloser, *smux.Config) (*smux.Session, error)
	if t.role == "source" {
		newSessFn = smux.Server
	} else {
		newSessFn = smux.Client
	}
	t.muxDataChannel = dc

	var conn io.ReadWriteCloser = pc1
	if !t.muxNocomp {
		conn = NewCompStream(conn)
	}
	t.muxSess, err = newSessFn(conn, t.getSmuxConfig())
	if err != nil {
		logger.WithError(err).Debugf("failed to upgrade smux conn")
		dc.Close()
		return
	}

	go t.muxSessionAcceptLoop()
	t.channelDone(1)
	logger.Tracef("on mux data channel open")
}

func (t *WebrtcTransport) onKcpDataChannelOpen(dc *webrtc.DataChannel) {
	var err error

	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":        "onKcpDataChannelOpen",
		"label":          dc.Label(),
		"smuxVer":        t.muxVer,
		"smuxBuf":        t.muxBuf,
		"smuxStreamBuf":  t.muxStreamBuf,
		"smuxKeepalive":  t.muxKeepalive,
		"smuxNocomp":     t.muxNocomp,
		"kcpPreset":      t.kcpPreset,
		"kcpCrypt":       t.kcpCrypt,
		"kcpMtu":         t.kcpMtu,
		"kcpSndwnd":      t.kcpSndwnd,
		"kcpRcvwnd":      t.kcpRcvwnd,
		"kcpDataShard":   t.kcpDataShard,
		"kcpParityShard": t.kcpParityShard,
	})

	t.kcpMtx.Lock()
	defer t.kcpMtx.Unlock()

	rwc, err := dc.Detach()
	if err != nil {
		logger.WithError(err).Debugf("failed to detach DataChannel")
		dc.Close()
		return
	}

	rwc1, err := t.upgradeKcpConn(rwc)
	if err != nil {
		logger.WithError(err).Debugf("failed to wrap kcp conn")
		dc.Close()
		return
	}

	var newSessFn func(io.ReadWriteCloser, *smux.Config) (*smux.Session, error)
	if t.role == "source" {
		newSessFn = smux.Server
	} else {
		newSessFn = smux.Client
	}
	t.kcpDataChannel = dc

	var conn io.ReadWriteCloser = rwc1
	if !t.muxNocomp {
		conn = NewCompStream(conn)
	}
	t.kcpSess, _ = newSessFn(conn, t.getSmuxConfig())

	go t.kcpSessionAcceptLoop()
	t.channelDone(1)
	logger.Tracef("on kcp data channel open")
}

func (t *WebrtcTransport) channelDone(n int32) {
	x := atomic.AddInt32(&t.readyCount, -n)
	if x <= 0 {
		t.readyOnce.Do(func() {
			t.onReadyCb(t)
			close(t.ready)
		})
	}
}