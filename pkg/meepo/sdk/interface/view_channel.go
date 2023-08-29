package sdk_interface

type ChannelView struct {
	Addr        string `json:"addr"`
	ID          uint16 `json:"id"`
	Mode        string `json:"mode"`
	State       string `json:"state"`
	IsSource    bool   `json:"isSource"`
	IsSink      bool   `json:"isSink"`
	SinkNetwork string `json:"sinkNetwork"`
	SinkAddress string `json:"sinkAddress"`
}
