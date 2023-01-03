package well_known_option

import "github.com/PeerXu/meepo/pkg/lib/option"

const OPTION_KCP_DATA_SHARD = "kcpDataShard"

var WithKcpDataShard, GetKcpDataShard = option.New[int](OPTION_KCP_DATA_SHARD)
