package teleportation

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/transport"
	mgroup "github.com/PeerXu/meepo/pkg/util/group"
)

type TeleportationSink struct {
	opt    ofn.Option
	logger logrus.FieldLogger

	name         string
	source       net.Addr
	sink         net.Addr
	transport    transport.Transport
	datachannels map[string]transport.DataChannel

	newDial NewDial

	onDoTeleportFunc OnDoTeleportFunc
	onCloseHandler   OnCloseHandler
	onErrorHandler   OnErrorHandler

	datachannelsMtx sync.Mutex
}

func (ts *TeleportationSink) getLogger() logrus.FieldLogger {
	return ts.logger.WithFields(logrus.Fields{
		"#instance": "TeleportationSink",
		"name":      ts.Name(),
	})
}

func (ts *TeleportationSink) onError(err error) {
	if ts.onErrorHandler != nil {
		ts.onErrorHandler(err)
	}

	ts.close()
}

func (ts *TeleportationSink) close() error {
	var eg errgroup.Group

	ts.datachannelsMtx.Lock()
	for _, dc := range ts.datachannels {
		eg.Go(dc.Close)
	}
	ts.datachannelsMtx.Unlock()

	return eg.Wait()
}

func (ts *TeleportationSink) OnDoTeleport(label string) error {
	logger := ts.getLogger().WithFields(logrus.Fields{
		"#method": "OnDoTeleport",
		"label":   label,
	})

	sink := ts.Sink()
	conn, err := ts.newDial(sink.Network(), sink.String())
	if err != nil {
		logger.WithError(err).Debugf("failed to dial to sink")
		return err
	}
	logger.Tracef("dial")

	tp := ts.Transport()
	tp.OnDataChannelCreate(label, func(dc transport.DataChannel) {
		rg := mgroup.NewRaceGroupFunc()

		outerDataChannelCloser := dc.Close
		defer func() {
			if outerDataChannelCloser != nil {
				outerDataChannelCloser()
			}
		}()

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
				innerLogger.WithError(dc.Close()).Tracef("datachannel closed")
				innerLogger.WithError(conn.Close()).Tracef("conn closed")

				ts.datachannelsMtx.Lock()
				delete(ts.datachannels, dc.Label())
				ts.datachannelsMtx.Unlock()
				innerLogger.Tracef("remove from data channels")

				innerLogger.WithError(err).Tracef("done")
			}()
			logger.Tracef("data channel opened")
		})

		ts.datachannelsMtx.Lock()
		ts.datachannels[label] = dc
		ts.datachannelsMtx.Unlock()
		logger.Tracef("add to data channels")

		outerDataChannelCloser = nil

		logger.Tracef("data channel created")
	})

	return nil
}

func (ts *TeleportationSink) Name() string {
	return ts.name
}

func (ts *TeleportationSink) Source() net.Addr {
	return ts.source
}

func (ts *TeleportationSink) Sink() net.Addr {
	return ts.sink
}

func (ts *TeleportationSink) Portal() Portal {
	return PortalSink
}

func (ts *TeleportationSink) Transport() transport.Transport {
	return ts.transport
}

func (ts *TeleportationSink) DataChannels() []transport.DataChannel {
	var dcs []transport.DataChannel

	ts.datachannelsMtx.Lock()
	for _, dc := range ts.datachannels {
		dcs = append(dcs, dc)
	}
	ts.datachannelsMtx.Unlock()

	return dcs
}

func (ts *TeleportationSink) Close() error {
	if ts.onCloseHandler != nil {
		ts.onCloseHandler()
	}

	return ts.close()
}

func newNewTeleportationSinkOptions() ofn.Option {
	return ofn.NewOption(map[string]interface{}{})
}

func NewTeleportationSink(opts ...NewTeleportationSinkOption) (*TeleportationSink, error) {
	var logger logrus.FieldLogger
	var name string
	var source, sink net.Addr
	var tp transport.Transport
	var newDial NewDial
	var odtf OnDoTeleportFunc
	var och OnCloseHandler
	var oeh OnErrorHandler
	var val interface{}
	var ok bool

	o := newNewTeleportationSinkOptions()

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

	if val = o.Get("newDial").Inter(); val != nil {
		newDial = val.(NewDial)
	} else {
		newDial = net.Dial
	}

	if val = o.Get("onDoTeleportFunc").Inter(); val != nil {
		odtf = val.(OnDoTeleportFunc)
	}

	if val = o.Get("onCloseHandler").Inter(); val != nil {
		och = val.(OnCloseHandler)
	}

	if val = o.Get("onErrorHandler").Inter(); val != nil {
		oeh = val.(OnErrorHandler)
	}

	ts := &TeleportationSink{
		opt:              o,
		logger:           logger,
		name:             name,
		source:           source,
		sink:             sink,
		transport:        tp,
		datachannels:     make(map[string]transport.DataChannel),
		newDial:          newDial,
		onDoTeleportFunc: odtf,
		onCloseHandler:   och,
		onErrorHandler:   oeh,
	}

	return ts, nil
}
