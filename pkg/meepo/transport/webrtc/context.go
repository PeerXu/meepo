package transport_webrtc

import "context"

func (t *WebrtcTransport) context() context.Context {
	return context.Background()
}
