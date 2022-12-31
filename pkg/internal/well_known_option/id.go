package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_ID = "id"
)

var (
	WithID, GetID = option.New[string](OPTION_ID)
)
