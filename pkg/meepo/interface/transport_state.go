package meepo_interface

type TransportState string

const (
	TRANSPORT_STATE_UNKNOWN      TransportState = "unknown"
	TRANSPORT_STATE_NEW          TransportState = "new"
	TRANSPORT_STATE_CONNECTING   TransportState = "connecting"
	TRANSPORT_STATE_CONNECTED    TransportState = "connected"
	TRANSPORT_STATE_DISCONNECTED TransportState = "disconnected"
	TRANSPORT_STATE_FAILED       TransportState = "failed"
	TRANSPORT_STATE_CLOSED       TransportState = "closed"
)

func (x TransportState) String() string {
	return string(x)
}
