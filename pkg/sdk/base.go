package sdk

import "net"

type BaseMeepoSDK struct {
	BaseTransportSDK
	BaseTeleportationSDK
}

func (BaseMeepoSDK) Version() (*Version, error) {
	return nil, UnimplementedError
}

func (BaseMeepoSDK) Ping() error {
	return UnimplementedError
}

func (BaseMeepoSDK) Shutdown() error {
	return UnimplementedError
}

func (BaseMeepoSDK) Whoami() (string, error) {
	return "", UnimplementedError
}

func (BaseMeepoSDK) Teleport(peerID string, remote net.Addr, opt *TeleportOption) (net.Addr, error) {
	return nil, UnimplementedError
}

type BaseTransportSDK struct{}

func (BaseTransportSDK) NewTransport(peerID string) (*Transport, error) {
	return nil, UnimplementedError
}

func (BaseTransportSDK) CloseTransport(peerID string) error {
	return UnimplementedError
}

func (BaseTransportSDK) ListTransports() ([]*Transport, error) {
	return nil, UnimplementedError
}

func (BaseTransportSDK) GetTransport(peerID string) (*Transport, error) {
	return nil, UnimplementedError
}

type BaseTeleportationSDK struct{}

func (BaseTeleportationSDK) NewTeleportation(peerID string, sink net.Addr, opt *NewTeleportationOption) (*Teleportation, error) {
	return nil, UnimplementedError
}

func (BaseTeleportationSDK) CloseTeleportation(name string) error {
	return UnimplementedError
}

func (BaseTeleportationSDK) ListTeleportations() ([]*Teleportation, error) {
	return nil, UnimplementedError
}

func (BaseTeleportationSDK) GetTeleportation(name string) (*Teleportation, error) {
	return nil, UnimplementedError
}
