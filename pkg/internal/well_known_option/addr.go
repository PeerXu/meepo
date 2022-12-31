package well_known_option

import (
	"github.com/PeerXu/meepo/pkg/internal/option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

const (
	OPTION_ADDR = "addr"
)

var (
	WithAddr, GetAddr = option.New[meepo_interface.Addr](OPTION_ADDR)
)
