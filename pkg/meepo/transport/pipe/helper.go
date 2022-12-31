package transport_pipe

import "sync/atomic"

func (pt *PipeTransport) nextChannelID() uint16 {
	return uint16(atomic.AddUint32(&pt.currentChannelID, 1) & 0xffff)
}
