package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_DESTINATION = "destination"
)

var (
	WithDestination, GetDestination = option.New[[]byte](OPTION_DESTINATION)
)
