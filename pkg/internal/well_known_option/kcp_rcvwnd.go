package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_KCP_RCVWND = "kcpRcvwnd"

var WithKcpRecvwnd, GetKcpRcvwnd = option.New[int](OPTION_KCP_RCVWND)
