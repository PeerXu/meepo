package transport_core

import "github.com/PeerXu/meepo/pkg/lib/option"

type TransportHooks struct {
	AfterNewTransportHook    AfterNewTransportHook
	BeforeCloseTransportHook BeforeCloseTransportHook
	AfterCloseTransportHook  AfterCloseTransportHook
}

func ApplyTransportHooks(o option.Option, h *TransportHooks) {
	h.AfterNewTransportHook, _ = GetAfterNewTransportHook(o)
	h.BeforeCloseTransportHook, _ = GetBeforeCloseTransportHook(o)
	h.AfterCloseTransportHook, _ = GetAfterCloseTransportHook(o)
}
