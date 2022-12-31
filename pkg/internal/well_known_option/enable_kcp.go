package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_ENABLE_KCP = "enableKcp"

var WithEnableKcp, GetEnableKcp = option.New[bool](OPTION_ENABLE_KCP)
