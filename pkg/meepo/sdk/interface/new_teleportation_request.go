package sdk_interface

type NewTeleportationRequest struct {
	Target string
	TeleportationParam
}

type TeleportationParam struct {
	SourceNetwork string
	SourceAddress string
	SinkNetwork   string
	SinkAddress   string
	Mode          string
}
