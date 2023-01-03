package well_known_option

import "github.com/PeerXu/meepo/pkg/lib/option"

const OPTION_ENABLE_MUX = "enableMux"

var WithEnableMux, GetEnableMux = option.New[bool](OPTION_ENABLE_MUX)
