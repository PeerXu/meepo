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
	mgroup "github.com/PeerXu/meepo/pkg/util/group"
)

type TeleportationSource struct {
	opt    objx.Map
	logger logrus.FieldLogger
	idx    int32
	closed int32

	name         string
	source       net.Addr
	sink         net.Addr
	transport    transport.Transport
	datachannels map[string]transport.DataChannel

	dialRequests chan *DialRequest

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

func (ts *TeleportationSource) onError(err error) {
	defer ts.close()

	if ts.onErrorHandler != nil {
		ts.onErrorHandler(err)
	}
}

func (ts *TeleportationSource) close() error {
	if atomic.SwapInt32(&ts.closed, 1) == 1 {
		return nil
	}

	var eg errgroup.Group

	ts.datachannelsMtx.Lock()
	for _, dc := range ts.datachannels {
		eg.Go(dc.Close)
	}
	ts.datachannelsMtx.Unlock()

	return eg.Wait()
}

func (ts *TeleportationSource) requestLoop() {
	logger := ts.getLogger().WithField("#method", "fetchConnLoop")

	defer logger.Tracef("exited")

	for {
		dr, ok := <-ts.dialRequests
		if !ok {
			return
		}
		go ts.onDial(dr)
	}
}

func (ts *TeleportationSource) onDial(dr *DialRequest) {
	var dc transport.DataChannel
	var err error

	conn := dr.Conn
	rg := mgroup.NewRaceGroupFunc()
	logger := ts.getLogger().WithField("#method", "onDial")

	outerConnCloser := conn.Close
	defer func() {
		if outerConnCloser != nil {
			outerConnCloser()
		}
	}()

	tp := ts.Transport()

	idx := atomic.AddInt32(&ts.idx, 1)
	label := fmt.Sprintf("%s:%d", ts.Name(), idx)

	if err := ts.doTeleportFunc(label); err != nil {
		logger.WithError(err).Debugf("failed to do teleport func")
		rg.Go(mgroup.DONE)
		return
	}
	logger.Tracef("do teleport func")

	dc, err = tp.CreateDataChannel(
		label,
		transport.WithOrdered(true),
	)
	if err != nil {
		logger.WithError(err).Debugf("failed to create data channel")
		rg.Go(mgroup.DONE)
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

		rg.Go(func() (interface{}, error) {
			_, err := io.Copy(dc, conn)
			innerLogger.WithError(err).Tracef("conn->dc closed")
			return nil, err
		})
		rg.Go(func() (interface{}, error) {
			_, err := io.Copy(conn, dc)
			innerLogger.WithError(err).Tracef("conn<-dc closed")
			return nil, err
		})
		go func() {
			_, err := rg.Wait()
			innerLogger.WithError(conn.Close()).Tracef("conn closed")
			innerLogger.WithError(dc.Close()).Tracef("datachannel closed")

			ts.datachannelsMtx.Lock()
			delete(ts.datachannels, dc.Label())
			ts.datachannelsMtx.Unlock()
			innerLogger.Tracef("remove from data channels")

			if dr.Quit != nil {
				close(dr.Quit)
				innerLogger.Tracef("send quit signal")
			}

			innerLogger.WithError(err).Tracef("done")
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
	var drc chan *DialRequest
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

	if val = o.Get("dialRequestChannel").Inter(); val != nil {
		drc = val.(chan *DialRequest)
	}

	ts := &TeleportationSource{
		opt:            o,
		logger:         logger,
		name:           name,
		source:         source,
		sink:           sink,
		transport:      tp,
		datachannels:   make(map[string]transport.DataChannel),
		doTeleportFunc: dtf,
		onCloseHandler: och,
		onErrorHandler: oeh,
		dialRequests:   drc,
	}

	go ts.requestLoop()

	return ts, nil
}
