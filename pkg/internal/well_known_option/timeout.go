package well_known_option

import (
	"time"

	"github.com/PeerXu/meepo/pkg/internal/option"
)

const (
	OPTION_TIMEOUT = "timeout"
)

var (
	WithTimeout, GetTimeout = option.New[time.Duration](OPTION_TIMEOUT)
)
