package well_known_option

import "github.com/PeerXu/meepo/pkg/lib/option"

const OPTION_MUX_BUF = "muxBuf"

var WithMuxBuf, GetMuxBuf = option.New[int](OPTION_MUX_BUF)
