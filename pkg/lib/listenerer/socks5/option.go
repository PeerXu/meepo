package listenerer_socks5

import (
	"time"

	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

const (
	NAME = "socks5"
)

func DefaultListenOption() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_CONN_WAIT_ENABLED_TIMEOUT: 121 * time.Second,
	})
}
