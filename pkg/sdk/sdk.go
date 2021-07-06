package sdk

import (
	"net"
	"sync"

	encoding_api "github.com/PeerXu/meepo/pkg/api/encoding"
	"github.com/PeerXu/meepo/pkg/ofn"
)

type Version = encoding_api.Version
type Transport = encoding_api.Transport
type Teleportation = encoding_api.Teleportation

type TeleportOption struct {
	Name   string
	Local  net.Addr
	Secret string
}

type MeepoSDK interface {
	TransportSDK
	TeleportationSDK

	Version() (*Version, error)
	Ping(peerID string) error
	Shutdown() error
	Whoami() (string, error)
	Teleport(peerID string, remote net.Addr, opt *TeleportOption) (net.Addr, error)
}

type TransportSDK interface {
	NewTransport(peerID string) (*Transport, error)
	CloseTransport(peerID string) error
	ListTransports() ([]*Transport, error)
	GetTransport(peerID string) (*Transport, error)
}

type NewTeleportationOption struct {
	Name   string
	Source net.Addr
	Secret string
}

type TeleportationSDK interface {
	NewTeleportation(peerID string, remote net.Addr, opt *NewTeleportationOption) (*Teleportation, error)
	CloseTeleportation(name string) error
	ListTeleportations() ([]*Teleportation, error)
	GetTeleportation(name string) (*Teleportation, error)
}

type NewMeepoSDKOption = ofn.OFN

type NewMeepoSDKFunc func(opts ...NewMeepoSDKOption) (MeepoSDK, error)

var (
	newMeepoSDKFuncs sync.Map
)

func RegisterNewMeepoSDKFunc(name string, fn NewMeepoSDKFunc) {
	newMeepoSDKFuncs.Store(name, fn)
}

func NewMeepoSDK(name string, opts ...NewMeepoSDKOption) (MeepoSDK, error) {
	fn, ok := newMeepoSDKFuncs.Load(name)
	if !ok {
		return nil, UnsupportedMeepoSDKDriverError(name)
	}

	return fn.(NewMeepoSDKFunc)(opts...)
}
