package simple_logger

import (
	logging "github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/lib/config"
)

func GetLogger() (logging.Logger, error) {
	cfg := config.Get()
	return logging.NewLogger(logging.WithLevel(cfg.Meepo.Log.Level))
}
