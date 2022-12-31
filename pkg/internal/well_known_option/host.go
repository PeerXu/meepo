package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_HOST = "host"
)

var (
	WithHost, GetHost = option.New[string](OPTION_HOST)
)
