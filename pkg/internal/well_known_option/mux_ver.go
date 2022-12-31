package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_MUX_VER = "muxVer"

var WithMuxVer, GetMuxVer = option.New[int](OPTION_MUX_VER)
