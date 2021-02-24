package transport

import "sync"

type TransportState int

const (
	TransportStateNew TransportState = iota + 1
	TransportStateConnecting
	TransportStateConnected
	TransportStateDisconnected
	TransportStateFailed
	TransportStateClosed
)

var (
	transportStateStr = []string{
		TransportStateNew:          "new",
		TransportStateConnecting:   "connecting",
		TransportStateConnected:    "connected",
		TransportStateDisconnected: "disconnected",
		TransportStateFailed:       "failed",
		TransportStateClosed:       "closed",
	}
)

func (t TransportState) String() string {
	return transportStateStr[t]
}

type HandleID = uint32
type OnTransportStateHandler func(HandleID)

type Transport interface {
	PeerID() string
	Err() error
	Close() error
	OnTransportStateChange(func(TransportState))
	OnTransportState(TransportState, func(hid HandleID)) HandleID
	UnsetOnTransportState(s TransportState, hid HandleID)
	TransportState() TransportState

	DataChannels() ([]DataChannel, error)
	DataChannel(label string) (DataChannel, error)
	CreateDataChannel(label string, opts ...CreateDataChannelOption) (DataChannel, error)
	OnDataChannelCreate(label string, f func(DataChannel))
}

var (
	newTransportFuncs sync.Map
)

type NewTransportFunc func(...NewTransportOption) (Transport, error)

func NewTransport(name string, opts ...NewTransportOption) (Transport, error) {
	fn, ok := newTransportFuncs.Load(name)
	if !ok {
		return nil, UnsupportedTransportError(name)
	}

	return fn.(NewTransportFunc)(opts...)
}

func RegisterNewTransportFunc(name string, fn NewTransportFunc) {
	newTransportFuncs.Store(name, fn)
}
