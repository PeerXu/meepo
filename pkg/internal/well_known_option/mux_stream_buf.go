package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_MUX_STREAM_BUF = "muxStreamBuf"

var WithMuxStreamBuf, GetMuxStreamBuf = option.New[int](OPTION_MUX_STREAM_BUF)
