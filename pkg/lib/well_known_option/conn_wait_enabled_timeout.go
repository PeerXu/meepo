package well_known_option

import (
	"time"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

const OPTION_CONN_WAIT_ENABLED_TIMEOUT = "connWaitEnabledTimeout"

var (
	WithConnWaitEnabledTimeout, GetConnWaitEnabledTimeoout = option.New[time.Duration](OPTION_CONN_WAIT_ENABLED_TIMEOUT)
)
