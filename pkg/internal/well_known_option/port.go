package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_PORT = "port"
)

var (
	WithPort, GetPort = option.New[string](OPTION_PORT)
)
