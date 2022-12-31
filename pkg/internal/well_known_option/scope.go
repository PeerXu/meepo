package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_SCOPE = "scope"
)

var (
	WithScope, GetScope = option.New[string](OPTION_SCOPE)
)
