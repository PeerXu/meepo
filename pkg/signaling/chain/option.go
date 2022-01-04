package chain_signaling

import (
	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/PeerXu/meepo/pkg/signaling"
)

func DefaultEngineOption() ofn.Option {
	return ofn.NewOption(map[string]interface{}{})
}

func WithEngine(engines ...signaling.Engine) signaling.NewEngineOption {
	return func(o ofn.Option) {
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
