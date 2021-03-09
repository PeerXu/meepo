package meepo

import "fmt"

type BroadcastGetter interface {
	GetBroadcast() *Broadcast
}

type Broadcast struct {
	SourceID         string `json:"_broadcast_sourceID"`
	DestinationID    string `json:"_broadcast_destinationID"`
	BroadcastSession int32  `json:"_broadcast_broadcastSession"`
	Hop              int32  `json:"_broadcast_hop"`
	DetectNextHop    bool   `json:"_broadcast_detectNextHop,omitempty"`
}

func (b *Broadcast) String() string {
	return fmt.Sprintf("#<Broadcast: SourceID: %v, DestinationID: %v, BroadcastSession: %v, Hop: %v, DetectNextHop: %v>", b.SourceID, b.DestinationID, b.BroadcastSession, b.Hop, b.DetectNextHop)
}

func (b *Broadcast) GetBroadcast() *Broadcast {
	return b
}

func (b *Broadcast) Identifier() string {
	return fmt.Sprintf("%v.%v", b.DestinationID, b.BroadcastSession)
}

type BroadcastResponse struct {
	*Message
	*Broadcast
}

func InvertBroadcast(b *Broadcast, id string) *Broadcast {
	return &Broadcast{
		SourceID:         id,
		DestinationID:    b.DestinationID,
		BroadcastSession: b.BroadcastSession,
		Hop:              b.Hop,
		DetectNextHop:    b.DetectNextHop,
	}
}

func NextHopBroadcast(b *Broadcast) *Broadcast {
	var x Broadcast
	x = *b
	if x.Hop > 0 {
		x.Hop -= 1
	}
	return &x
}
