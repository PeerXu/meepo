package transport_webrtc

import (
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) onSystemRequestMessage(in Message) {
	logger := t.GetLogger().WithField("#method", "onSystemRequestMessage").WithFields(t.wrapMessage(in))

	ctx := marshaler.ContextWithMarshalerAndUnmarshaler(t.context(), t.marshaler, t.unmarshaler)
	var h meepo_interface.HandleFunc
	switch in.Method {
	case SYS_METHOD_PING:
		h = transport_core.WrapHandleFuncGenerics(t.onPing)
	case SYS_METHOD_NEW_CHANNEL:
		h = transport_core.WrapHandleFuncGenerics(t.onNewChannel)
	case SYS_METHOD_ADD_PEER_CONNECTION:
		h = transport_core.WrapHandleFuncGenerics(t.onAddPeerConnection)
	case SYS_METHOD_CLOSE:
		h = transport_core.WrapHandleFuncGenerics(t.onClose)
	default:
		var ok bool
		method := "sys/" + in.Method
		t.fnsMtx.Lock()
		h, ok = t.fns[method]
		t.fnsMtx.Unlock()
		if !ok {
			err := ErrUnsupportedMethodFn(method)
			t.sendErrorResponse(ctx, in, err)
			logger.WithError(err).Debugf("unsupported method")
			return
		}
	}

	buf, err := h(ctx, in.Data)
	if err != nil {
		t.sendErrorResponse(ctx, in, err)
		logger.WithError(err).Debugf("failed to handle system request")
		return
	}

	t.sendResponse(ctx, in, buf)

	logger.Tracef("on system request message")
}
