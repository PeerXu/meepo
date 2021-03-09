package chain_signaling

import (
	"github.com/PeerXu/meepo/pkg/signaling"
	"github.com/stretchr/objx"
)

func DefaultEngineOption() objx.Map {
	return objx.New(map[string]interface{}{})
}

func WithEngine(engines ...signaling.Engine) signaling.NewEngineOption {
	return func(o objx.Map) {
		var enginesSlice []signaling.Engine
		var ok bool

		enginesSlice, ok = o.Get("engines").Inter().([]signaling.Engine)
		if ok {
			enginesSlice = append(enginesSlice, engines...)
		} else {
			enginesSlice = engines
		}

		o["engines"] = enginesSlice
	}
}
