package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_KCP_PRESET = "kcpPreset"

var WithKcpPreset, GetKcpPreset = option.New[string](OPTION_KCP_PRESET)
