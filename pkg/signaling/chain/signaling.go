package chain_signaling

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/objx"

	"github.com/PeerXu/meepo/pkg/signaling"
)

type ChainEngine struct {
	opt    objx.Map
	logger logrus.FieldLogger

	engines []signaling.Engine
}

func (e *ChainEngine) Wire(dst, src *signaling.Descriptor) (*signaling.Descriptor, error) {
	var res *signaling.Descriptor
	var err error

	for _, ng := range e.engines {
		if res, err = ng.Wire(dst, src); err == nil {
			return res, nil
		}
	}

	return nil, err
}

func (e *ChainEngine) OnWire(handler signaling.WireHandler) {
	for _, ng := range e.engines {
		ng.OnWire(handler)
	}
}

func (e *ChainEngine) Close() error {
	var errOnce sync.Once
	var err error

	for _, ng := range e.engines {
		if er := ng.Close(); er != nil {
			errOnce.Do(func() { err = er })
		}
	}

	return err
}

func NewChainEngine(opts ...signaling.NewEngineOption) (signaling.Engine, error) {
	o := DefaultEngineOption()

	for _, opt := range opts {
		opt(o)
	}

	logger, ok := o.Get("logger").Inter().(logrus.FieldLogger)
	if !ok {
		return nil, fmt.Errorf("Require logger")
	}

	engines, ok := o.Get("engines").Inter().([]signaling.Engine)
	if !ok {
		return nil, fmt.Errorf("Require engines")
	}

	ce := &ChainEngine{
		opt:     o,
		logger:  logger,
		engines: engines,
	}

	return ce, nil
}

func init() {
	signaling.RegisterNewEngineFunc("chain", NewChainEngine)
}
