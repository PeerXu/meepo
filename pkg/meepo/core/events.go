package meepo_core

import (
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

const (
	EVENT_TRANSPORT_ACTION_NEW   = "mpo.transport.action.new"
	EVENT_TRANSPORT_ACTION_CLOSE = "mpo.transport.action.close"

	EVENT_CHANNEL_ACTION_NEW   = "mpo.channel.action.new"
	EVENT_CHANNEL_ACTION_CLOSE = "mpo.channel.action.close"

	EVENT_TELEPORTATION_ACTION_NEW   = "mpo.teleportation.action.new"
	EVENT_TELEPORTATION_ACTION_CLOSE = "mpo.teleportation.action.close"

	EVENT_TRANSPORT_STATE_NEW          = "mpo.transport.state.new"
	EVENT_TRANSPORT_STATE_CONNECTING   = "mpo.transport.state.connecting"
	EVENT_TRANSPORT_STATE_CONNECTED    = "mpo.transport.state.connected"
	EVENT_TRANSPORT_STATE_DISCONNECTED = "mpo.transport.state.disconnected"
	EVENT_TRANSPORT_STATE_FAILED       = "mpo.transport.state.failed"
	EVENT_TRANSPORT_STATE_CLOSED       = "mpo.transport.state.closed"

	EVENT_CHANNEL_STATE_NEW        = "mpo.channel.state.new"
	EVENT_CHANNEL_STATE_CONNECTING = "mpo.channel.state.connecting"
	EVENT_CHANNEL_STATE_OPEN       = "mpo.channel.state.open"
	EVENT_CHANNEL_STATE_CLOSING    = "mpo.channel.state.closing"
	EVENT_CHANNEL_STATE_CLOSED     = "mpo.channel.state.closed"
)

var (
	transportStateMap = map[meepo_interface.TransportState]string{
		meepo_interface.TRANSPORT_STATE_NEW:          EVENT_TRANSPORT_STATE_NEW,
		meepo_interface.TRANSPORT_STATE_CONNECTING:   EVENT_TRANSPORT_STATE_CONNECTING,
		meepo_interface.TRANSPORT_STATE_CONNECTED:    EVENT_TRANSPORT_STATE_CONNECTED,
		meepo_interface.TRANSPORT_STATE_DISCONNECTED: EVENT_TRANSPORT_STATE_DISCONNECTED,
		meepo_interface.TRANSPORT_STATE_FAILED:       EVENT_TRANSPORT_STATE_FAILED,
		meepo_interface.TRANSPORT_STATE_CLOSED:       EVENT_TRANSPORT_STATE_CLOSED,
	}

	channelStateMap = map[meepo_interface.ChannelState]string{
		meepo_interface.CHANNEL_STATE_NEW:        EVENT_CHANNEL_STATE_NEW,
		meepo_interface.CHANNEL_STATE_CONNECTING: EVENT_CHANNEL_STATE_CONNECTING,
		meepo_interface.CHANNEL_STATE_OPEN:       EVENT_CHANNEL_STATE_OPEN,
		meepo_interface.CHANNEL_STATE_CLOSING:    EVENT_CHANNEL_STATE_CLOSING,
		meepo_interface.CHANNEL_STATE_CLOSED:     EVENT_CHANNEL_STATE_CLOSED,
	}
)

func ConvertTransportStateToEventName(st meepo_interface.TransportState) string {
	return transportStateMap[st]
}

func ConvertChannelStateToEventName(st meepo_interface.ChannelState) string {
	return channelStateMap[st]
}
