package well_known_option

import "github.com/PeerXu/meepo/pkg/lib/option"

const OPTION_ENABLE_KCP = "enableKcp"

var WithEnableKcp, GetEnableKcp = option.New[bool](OPTION_ENABLE_KCP)
