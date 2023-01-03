package transport_core

import (
	"time"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	OPTION_TRANSPORT               = "transport"
	OPTION_ON_TRANSPORT_CLOSE      = "onTransportClose"
	OPTION_ON_TRANSPORT_READY      = "onTransportReady"
	OPTION_BEFORE_NEW_CHANNEL_HOOK = "beforeNewChannelHook"
	OPTION_READY_TIMEOUT           = "readyTimeout"
)

type OnTransportCloseFunc = func(Transport) error

type OnTransportReadyFunc = func(Transport) error

type OnChannelCloseFunc = func(Channel) error

type BeforeNewChannelHook = func(t Transport, network, address string) error

type NewTransportOption = option.ApplyOption

var (
	WithOnTransportCloseFunc, GetOnTransportCloseFunc = option.New[OnTransportCloseFunc](OPTION_ON_TRANSPORT_CLOSE)
	WithOnTransportReadyFunc, GetOnTransportReadyFunc = option.New[OnTransportReadyFunc](OPTION_ON_TRANSPORT_READY)
	WithBeforeNewChannelHook, GetBeforeNewChannelHook = option.New[BeforeNewChannelHook](OPTION_BEFORE_NEW_CHANNEL_HOOK)
	WithTransport, GetTransport                       = option.New[Transport](OPTION_TRANSPORT)
	WithReadyTimeout, GetReadyTimeout                 = option.New[time.Duration](OPTION_READY_TIMEOUT)
)
