package constant

const (
	WEBRTC_RECEIVE_BUFFER_SIZE uint32 = 33554432

	TRACKER_ADDR = "62vv3lwalqmdb2657f7ax73fem7gkgzmin3w7qyy0sjjfae0f3p"
	TRACKER_HOST = "tkd.mpo.peerstud.io:12346"

	API_HOST = "127.0.0.1:12345"

	SOCKS5_HOST = "127.0.0.1:12341"

	ACL = `- allow: "*"`

	LOG_LEVEL = "info"

	SMUX_VERSION            int  = 2
	SMUX_BUFFER_SIZE        int  = 4194304
	SMUX_STREAM_BUFFER_SIZE int  = 2097152
	SMUX_NOCOMP             bool = false

	KCP_PRESET       string = "normal"
	KCP_CRYPT        string = "null"
	KCP_KEY          string = "1m4s3cr3t"
	KCP_MTU          int    = 32767
	KCP_SNDWND       int    = 128
	KCP_RCVWND       int    = 512
	KCP_DATA_SHARD   int    = 10
	KCP_PARITY_SHARD int    = 10
)