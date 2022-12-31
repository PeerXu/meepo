package tracker_interface

import "github.com/pion/webrtc/v3"

type NewTransportRequest struct {
	Offer webrtc.SessionDescription

	EnableMux    bool
	MuxLabel     string
	MuxVer       int
	MuxBuf       int
	MuxStreamBuf int
	MuxNocomp    bool

	EnableKcp   bool
	KcpLabel    string
	KcpPreset   string
	KcpCrypt    string
	KcpKey      string
	KcpMtu      int
	KcpSndwnd   int
	KcpRcvwnd   int
	DataShard   int
	ParityShard int
}
