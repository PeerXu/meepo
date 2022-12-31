package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_MODE = "mode"

var WithMode, GetMode = option.New[string](OPTION_MODE)
