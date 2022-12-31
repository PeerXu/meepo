package listenerer

import (
	listenerer_interface "github.com/PeerXu/meepo/pkg/internal/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/internal/option"
)

const (
	OPTION_LISTENER   = "listener"
	OPTION_LISTENERER = "listenerer"
)

var (
	WithListener, GetListener = option.New[listenerer_interface.Listener](OPTION_LISTENER)
)
