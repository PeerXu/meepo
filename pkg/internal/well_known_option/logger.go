package well_known_option

import (
	"github.com/sirupsen/logrus"

	"github.com/PeerXu/meepo/pkg/internal/option"
)

const (
	OPTION_LOGGER = "logger"
)

var (
	WithLogger, GetLogger = option.New[*logrus.Entry](OPTION_LOGGER)
)
