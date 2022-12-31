package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_KCP_PARITY_SHARD = "kcpParityShard"

var WithKcpParityShard, GetKcpParityShard = option.New[int](OPTION_KCP_PARITY_SHARD)
