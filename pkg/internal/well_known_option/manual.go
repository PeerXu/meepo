package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_MANUAL = "manual"

var WithManual, GetManual = option.New[bool](OPTION_MANUAL)
