package sdk_core

import "github.com/PeerXu/meepo/pkg/lib/option"

const (
	METHOD_GET_VERSION             = "getVersion"
	METHOD_WHOAMI                  = "whoami"
	METHOD_PING                    = "ping"
	METHOD_DIAGNOSTIC              = "diagnostic"
	METHOD_NEW_TELEPORTATION       = "newTeleportation"
	METHOD_CLOSE_TELEPORTATION     = "closeTeleportation"
	METHOD_GET_TELEPORTATION       = "getTeleportation"
	METHOD_LIST_TELEPORTATIONS     = "listTeleportations"
	METHOD_TELEPORT                = "teleport"
	METHOD_NEW_TRANSPORT           = "newTransport"
	METHOD_CLOSE_TRANSPORT         = "closeTransport"
	METHOD_GET_TRANSPORT           = "getTransport"
	METHOD_LIST_TRANSPORTS         = "listTransports"
	METHOD_CLOSE_CHANNEL           = "closeChannel"
	METHOD_GET_CHANNEL             = "getChannel"
	METHOD_LIST_CHANNELS           = "listChannels"
	METHOD_LIST_CHANNELS_BY_TARGET = "listChannelsByTarget"
)

type NewSDKOption = option.ApplyOption
