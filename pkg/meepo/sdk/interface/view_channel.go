package sdk_interface

type ChannelView struct {
	Addr        string
	ID          uint16
	Mode        string
	State       string
	IsSource    bool
	IsSink      bool
	SinkNetwork string
	SinkAddress string
}
