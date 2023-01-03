package listenerer

import (
	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	OPTION_LISTENER   = "listener"
	OPTION_LISTENERER = "listenerer"
)

var (
	WithListener, GetListener = option.New[listenerer_interface.Listener](OPTION_LISTENER)
)
