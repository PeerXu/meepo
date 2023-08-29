package sdk_interface

type TeleportationView struct {
	ID            string `json:"id"`
	Mode          string `json:"mode"`
	Addr          string `json:"addr"`
	SourceNetwork string `json:"sourceNetwork"`
	SourceAddress string `json:"sourceAddress"`
	SinkNetwork   string `json:"sinkNetwork"`
	SinkAddress   string `json:"sinkAddress"`
}
