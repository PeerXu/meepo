package transport_pipe

import (
	"time"

	"github.com/PeerXu/meepo/pkg/internal/option"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func defaultNewPipeTransportOptions() option.Option {
	return option.NewOption(map[string]any{
		transport_core.OPTION_READY_TIMEOUT: 601 * time.Second,
	})
}
