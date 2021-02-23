package teleportation

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"
	"golang.org/x/sync/errgroup"

	"github.com/PeerXu/meepo/pkg/transport"
	ieg "github.com/PeerXu/meepo/pkg/util/errgroup"
)

type TeleportationSource struct {
	opt    objx.Map
	logger logrus.FieldLogger
	idx    int64

	name         string
	source       net.Addr
	sink         net.Addr
	transport    transport.Transport
	datachannels map[string]transport.DataChannel

	newListener NewListener
	lis         net.Listener

	doTeleportFunc DoTeleportFunc
	onCloseHandler OnCloseHandler
	onErrorHandler OnErrorHandler

	datachannelsMtx sync.Mutex
	lisMtx          sync.Mutex
}

func (ts *TeleportationSource) getLogger() logrus.FieldLogger {
	return ts.logger.WithFields(logrus.Fields{
		"#instance": "TeleportationSource",
		"name":      ts.Name(),
	})
}

func (ts *TeleportationSource) getOrCreateListener() (net.Listener, error) {
	var err error

	ts.lisMtx.Lock()
	defer ts.lisMtx.Unlock()

	if ts.lis == nil {
		if ts.lis, err = ts.newListener(ts.source.Network(), ts.source.String()); err != nil {
			return nil, err
		}
	}

	return ts.lis, nil
}

func (ts *TeleportationSource) onError(err error) {
	defer ts.close()

	if ts.onErrorHandler != nil {
		ts.onErrorHandler(err)
	}
}

func (ts *TeleportationSource) close() error {
	var eg errgroup.Group

	ts.lisMtx.Lock()
	if ts.lis != nil {
		eg.Go(ts.lis.Close)
	}
	ts.lisMtx.Unlock()

	ts.datachannelsMtx.Lock()
	for _, dc := range ts.datachannels {
		eg.Go(dc.Close)
	}
	ts.datachannelsMtx.Unlock()

	return eg.Wait()
}

func (ts *TeleportationSource) acceptLoop() {
	logger := ts.getLogger().WithField("#method", "acceptLoop")

	var conn net.Conn
	var err error

	defer logger.Tracef("accept loop exited")

	lis, err := ts.getOrCreateListener()
	if err != nil {
		logger.WithError(err).Debugf("failed to new listener")
		return
	}

	for {

		if conn, err = lis.Accept(); err != nil {
			if ts.onErrorHandler != nil && err != io.EOF {
				ts.onError(err)
			}
			return
		}

		go ts.onAccept(conn)
	}
}

func (ts *TeleportationSource) onAccept(conn net.Conn) {
	var wg sync.WaitGroup
	var eg ieg.ImmediatelyErrorGroup
	var dc transport.DataChannel
	var err error

	logger := ts.getLogger().WithField("#method", "onAccept")

	defer func() {
		wg.Wait()
	}()

	outerConnCloser := conn.Close
	defer func() {
		if outerConnCloser != nil {
			outerConnCloser()
		}
	}()

	tp := ts.Transport()

	idx := atomic.AddInt64(&ts.idx, 1)
	label := fmt.Sprintf("%s:%d", ts.Name(), idx)

	if err := ts.doTeleportFunc(label); err != nil {
		logger.WithError(err).Debugf("failed to do teleport func")
		eg.Go(func() error { return nil })
		return
	}
	logger.Tracef("do teleport func")

	dc, err = tp.CreateDataChannel(
		label,
		transport.WithOrdered(true),
	)
	if err != nil {
		logger.WithError(err).Debugf("failed to create data channel")
		eg.Go(func() error { return nil })
		return
	}
	outerDataChannelCloser := dc.Close
	defer func() {
		if outerDataChannelCloser != nil {
			outerDataChannelCloser()
		}
	}()
	logger.Tracef("create data channel")

	dc.OnOpen(func() {
		innerLogger := ts.getLogger().WithFields(logrus.Fields{
			"#method": "dataChannelCopyLoop",
			"label":   dc.Label(),
		})

		eg.Go(func() error {
			_, err := io.Copy(dc, conn)
			innerLogger.WithError(err).Tracef("conn->dc closed")
			return err
		})
		eg.Go(func() error {
			_, err := io.Copy(conn, dc)
			innerLogger.WithError(err).Tracef("conn<-dc closed")
			return err
		})
		go func() {
			innerLogger.WithError(eg.Wait()).Tracef("broken")
			innerLogger.WithError(conn.Close()).Tracef("conn closed")
			innerLogger.WithError(dc.Close()).Tracef("datachannel closed")

			ts.datachannelsMtx.Lock()
			delete(ts.datachannels, dc.Label())
			ts.datachannelsMtx.Unlock()
			innerLogger.Tracef("remove from data channels")

			innerLogger.Tracef("done")
		}()

		logger.Tracef("data channel opened")
	})

	ts.datachannelsMtx.Lock()
	ts.datachannels[label] = dc
	ts.datachannelsMtx.Unlock()
	logger.Tracef("add to data channels")

	outerConnCloser = nil
	outerDataChannelCloser = nil

	logger.Tracef("done")
}

func (ts *TeleportationSource) Name() string {
	return ts.name
}

func (ts *TeleportationSource) Source() net.Addr {
	return ts.source
}

func (ts *TeleportationSource) Sink() net.Addr {
	return ts.sink
}

func (ts *TeleportationSource) Portal() Portal {
	return PortalSource
}

func (ts *TeleportationSource) Transport() transport.Transport {
	return ts.transport
}

func (ts *TeleportationSource) DataChannels() []transport.DataChannel {
	var dcs []transport.DataChannel

	ts.datachannelsMtx.Lock()
	for _, dc := range ts.datachannels {
		dcs = append(dcs, dc)
	}
	ts.datachannelsMtx.Unlock()

	return dcs
}

func (ts *TeleportationSource) Close() error {
	if ts.onCloseHandler != nil {
		ts.onCloseHandler()
	}

	return ts.close()
}

func newNewteleportationSourceOptions() objx.Map {
	return objx.New(map[string]interface{}{})
}

func NewTeleportationSource(opts ...NewTeleportationSourceOption) (*TeleportationSource, error) {
	var logger logrus.FieldLogger
	var name string
	var source, sink net.Addr
	var tp transport.Transport
	var ok bool
	var dtf DoTeleportFunc
	var och OnCloseHandler
	var oeh OnErrorHandler
	var nl NewListener
	var val interface{}

	o := newNewteleportationSourceOptions()

	for _, opt := range opts {
		opt(o)
	}

	if logger, ok = o.Get("logger").Inter().(logrus.FieldLogger); !ok {
		logger = logrus.New()
	}

	if source, ok = o.Get("source").Inter().(net.Addr); !ok {
		return nil, fmt.Errorf("Require source")
	}

	if sink, ok = o.Get("sink").Inter().(net.Addr); !ok {
		return nil, fmt.Errorf("Require sink")
	}

	if tp, ok = o.Get("transport").Inter().(transport.Transport); !ok {
		return nil, fmt.Errorf("Require transport")
	}

	if name, ok = o.Get("name").Inter().(string); !ok {
		name = fmt.Sprintf("%s:%s", sink.Network(), sink.String())
	}

	if val = o.Get("doTeleportFunc").Inter(); val != nil {
		dtf = val.(DoTeleportFunc)
	}

	if val = o.Get("onCloseHandler").Inter(); val != nil {
		och = val.(OnCloseHandler)
	}

	if val = o.Get("onErrorHandler").Inter(); val != nil {
		oeh = val.(OnErrorHandler)
	}

	if nl, ok = o.Get("newListener").Inter().(NewListener); !ok {
		nl = net.Listen
	}

	ts := &TeleportationSource{
		opt:            o,
		logger:         logger,
		name:           name,
		source:         source,
		sink:           sink,
		transport:      tp,
		datachannels:   make(map[string]transport.DataChannel),
		newListener:    nl,
		doTeleportFunc: dtf,
		onCloseHandler: och,
		onErrorHandler: oeh,
	}

	go ts.acceptLoop()

	return ts, nil
}
