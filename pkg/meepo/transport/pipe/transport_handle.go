package transport_pipe

import (
	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *PipeTransport) Handle(method string, fn meepo_interface.HandleFunc, opts ...meepo_interface.HandleOption) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "Handle",
		"method":  method,
	})

	t.fnsMtx.Lock()
	defer t.fnsMtx.Unlock()

	t.fns[method] = fn

	logger.Tracef("handle method")
}
