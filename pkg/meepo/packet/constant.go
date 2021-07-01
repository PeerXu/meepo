package packet

type Type string

const (
	Request           Type = "request"
	Response          Type = "response"
	BroadcastRequest  Type = "broadcastRequest"
	BroadcastResponse Type = "broadcastResponse"
)
