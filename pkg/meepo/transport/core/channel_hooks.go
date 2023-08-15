package transport_core

import "github.com/PeerXu/meepo/pkg/lib/option"

type ChannelHooks struct {
	BeforeNewChannelHook   BeforeNewChannelHook
	AfterNewChannelHook    AfterNewChannelHook
	BeforeCloseChannelHook BeforeCloseChannelHook
	AfterCloseChannelHook  AfterCloseChannelHook
}

func ApplyChannelHooks(o option.Option, h *ChannelHooks) {
	h.BeforeNewChannelHook, _ = GetBeforeNewChannelHook(o)
	h.AfterNewChannelHook, _ = GetAfterNewChannelHook(o)
	h.BeforeCloseChannelHook, _ = GetBeforeCloseChannelHook(o)
	h.AfterCloseChannelHook, _ = GetAfterCloseChannelHook(o)
}
