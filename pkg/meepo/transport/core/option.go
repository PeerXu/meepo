package transport_core

import (
	"time"

	"github.com/PeerXu/meepo/pkg/lib/option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

const (
	OPTION_TRANSPORT                   = "transport"
	OPTION_ON_TRANSPORT_READY          = "onTransportReady"
	OPTION_BEFORE_NEW_TRANSPORT_HOOK   = "beforeNewTransportHook"
	OPTION_AFTER_NEW_TRANSPORT_HOOK    = "afterNewTransportHook"
	OPTION_BEFORE_CLOSE_TRANSPORT_HOOK = "beforeCloseTransportHook"
	OPTION_AFTER_CLOSE_TRANSPORT_HOOK  = "afterCloseTransportHook"
	OPTION_BEFORE_NEW_CHANNEL_HOOK     = "beforeNewChannelHook"
	OPTION_AFTER_NEW_CHANNEL_HOOK      = "afterNewChannelHook"
	OPTION_BEFORE_CLOSE_CHANNEL_HOOK   = "beforeCloseChannelHook"
	OPTION_AFTER_CLOSE_CHANNEL_HOOK    = "afterCloseChannelHook"
	OPTION_READY_TIMEOUT               = "readyTimeout"
	OPTION_IS_SOURCE                   = "isSource"
	OPTION_IS_SINK                     = "isSink"
)

type OnTransportReadyFunc = func(Transport) error

type OnChannelCloseFunc = func(Channel) error

type BeforeNewTransportHook = func(meepo_interface.Addr, ...HookOption) error
type AfterNewTransportHook = func(Transport, ...HookOption)
type BeforeCloseTransportHook = func(Transport, ...HookOption) error
type AfterCloseTransportHook = func(Transport, ...HookOption)

type BeforeNewChannelHook = func(network, address string, opts ...HookOption) error
type AfterNewChannelHook = func(c Channel, opts ...HookOption)
type BeforeCloseChannelHook = func(c Channel, opts ...HookOption) error
type AfterCloseChannelHook = func(c Channel, opts ...HookOption)

type NewTransportOption = option.ApplyOption

var (
	WithOnTransportReadyFunc, GetOnTransportReadyFunc = option.New[OnTransportReadyFunc](OPTION_ON_TRANSPORT_READY)

	WithBeforeNewTransportHook, GetBeforeNewTransportHook     = option.New[BeforeNewTransportHook](OPTION_BEFORE_NEW_TRANSPORT_HOOK)
	WithAfterNewTransportHook, GetAfterNewTransportHook       = option.New[AfterNewTransportHook](OPTION_AFTER_NEW_TRANSPORT_HOOK)
	WithBeforeCloseTransportHook, GetBeforeCloseTransportHook = option.New[BeforeCloseTransportHook](OPTION_BEFORE_CLOSE_TRANSPORT_HOOK)
	WithAfterCloseTransportHook, GetAfterCloseTransportHook   = option.New[AfterCloseTransportHook](OPTION_AFTER_CLOSE_TRANSPORT_HOOK)

	WithBeforeNewChannelHook, GetBeforeNewChannelHook     = option.New[BeforeNewChannelHook](OPTION_BEFORE_NEW_CHANNEL_HOOK)
	WithAfterNewChannelHook, GetAfterNewChannelHook       = option.New[AfterNewChannelHook](OPTION_AFTER_NEW_CHANNEL_HOOK)
	WithBeforeCloseChannelHook, GetBeforeCloseChannelHook = option.New[BeforeCloseChannelHook](OPTION_BEFORE_CLOSE_CHANNEL_HOOK)
	WithAfterCloseChannelHook, GetAfterCloseChannelHook   = option.New[AfterCloseChannelHook](OPTION_AFTER_CLOSE_CHANNEL_HOOK)

	WithTransport, GetTransport       = option.New[Transport](OPTION_TRANSPORT)
	WithReadyTimeout, GetReadyTimeout = option.New[time.Duration](OPTION_READY_TIMEOUT)
)

type HookOption = option.ApplyOption

var (
	WithIsSource, GetIsSource = option.New[bool](OPTION_IS_SOURCE)
	WithIsSink, GetIsSink     = option.New[bool](OPTION_IS_SINK)
)
