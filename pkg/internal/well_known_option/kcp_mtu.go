package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_KCP_MTU = "kcpMtu"

var WithKcpMtu, GetKcpMtu = option.New[int](OPTION_KCP_MTU)
