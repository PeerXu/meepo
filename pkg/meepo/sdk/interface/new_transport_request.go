package sdk_interface

type NewTransportRequest struct {
	Target string
	Manual bool
	SmuxParam
	KcpParam
}

type SmuxParam struct {
	EnableMux    bool
	MuxVer       int
	MuxBuf       int
	MuxStreamBuf int
	MuxNocomp    bool
}

type KcpParam struct {
	EnableKcp      bool
	KcpPreset      string
	KcpKey         string
	KcpCrypt       string
	KcpMtu         int
	KcpSndwnd      int
	KcpRcvwnd      int
	KcpDataShard   int
	KcpParityShard int
}
