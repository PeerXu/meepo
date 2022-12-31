package well_known_option

import (
	"math/rand"

	"github.com/PeerXu/meepo/pkg/internal/option"
)

const (
	OPTION_RAND_SOURCE = "randSource"
)

var (
	WithRandSource, GetRandSource = option.New[rand.Source](OPTION_RAND_SOURCE)
)
