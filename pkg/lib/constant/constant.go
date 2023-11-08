package constant

import "time"

var (
	WEBRTC_RECEIVE_BUFFER_SIZE uint32 = 33554432

	POOF_INTERVAL           = 31 * time.Second
	POOF_REQUEST_CANDIDATES = 3

	TRACKER_ADDR = "62vv3lwalqmdb2657f7ax73fem7gkgzmin3w7qyy0sjjfae0f3p"
	TRACKER_HOST = "tkd-0.meepo.dev:12346"

	TRACKERD_HOST = "127.0.0.1:12346"

	API_HOST = "127.0.0.1:12345"

	SOCKS5_HOST = "127.0.0.1:12341"

	ACL_REPLACEME = `REPLACE_ME`
	ACL_BLOCK_ALL = `- block: "*"`
	ACL_ALLOW_ALL = `- allow: "*"`

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

	PROFILE_MAIN  = "main"
	PROFILE_MINOR = "minor"
	PROFILE_DEV   = "dev"
)
