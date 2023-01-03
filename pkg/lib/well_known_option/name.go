package well_known_option

import "github.com/PeerXu/meepo/pkg/lib/option"

const (
	OPTION_NAME = "name"
)

var (
	WithName, GetName = option.New[string](OPTION_NAME)
)
