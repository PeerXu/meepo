package transport_webrtc

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) onHandle(ctx context.Context, method string, in meepo_interface.HandleRequest) (out meepo_interface.HandleResponse, err error) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "onHandle",
		"method":  method,
	})

	t.fnsMtx.Lock()
	fn, ok := t.fns[method]
	t.fnsMtx.Unlock()

	if !ok {
		err = transport_core.ErrUnsupportedMethodFn(method)
		logger.WithError(err).Debugf("unsupported method")
		return nil, err
	}

	out, err = fn(t.wrapCallContext(ctx), in)
	if err != nil {
		logger.WithError(err).Debugf("failed to call")
		return nil, err
	}

	logger.Tracef("on handle")

	return out, nil
}

func (t *WebrtcTransport) wrapCallContext(ctx context.Context) context.Context {
	nctx := marshaler.ContextWithMarshalerAndUnmarshaler(ctx, t.marshaler, t.unmarshaler)
	nctx = transport_core.ContextWithTransport(nctx, t)
	return nctx
}
