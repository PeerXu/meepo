package transport_pipe

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *PipeTransport) Call(ctx context.Context, method string, req meepo_interface.CallRequest, res meepo_interface.CallResponse, opts ...meepo_interface.CallOption) (err error) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "Call",
		"method":  method,
	})

	t.fnsMtx.Lock()
	fn, ok := t.fns[method]
	t.fnsMtx.Unlock()

	if !ok {
		err = transport_core.ErrUnsupportedMethodFn(method)
		logger.WithError(err).Debugf("unsupported method")
		return err
	}

	in, err := t.marshaler.Marshal(req)
	if err != nil {
		logger.WithError(err).Debugf("failed to marshal request")
		return err
	}

	out, err := fn(t.wrapCallContext(ctx), in)
	if err != nil {
		logger.WithError(err).Debugf("failed to call")
		return err
	}

	if err = t.unmarshaler.Unmarshal(out, res); err != nil {
		logger.WithError(err).Debugf("failed to unmarshal response")
		return err
	}

	logger.Tracef("call")

	return nil
}

func (t *PipeTransport) wrapCallContext(ctx context.Context) context.Context {
	nctx := marshaler.ContextWithMarshalerAndUnmarshaler(ctx, t.marshaler, t.unmarshaler)
	nctx = transport_core.ContextWithTransport(nctx, t)
	return nctx
}
