package simple_logger

import (
	"github.com/PeerXu/meepo/pkg/lib/config"
	logging "github.com/PeerXu/meepo/pkg/lib/logging"
)

func GetLogger() (logging.Logger, error) {
	cfg := config.Get()

	return logging.NewLogger(
		logging.WithLevel(cfg.Meepo.Log.Level),
		logging.WithFile(cfg.Meepo.Log.File),
	)
}
