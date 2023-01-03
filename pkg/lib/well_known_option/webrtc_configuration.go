package well_known_option

import (
	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	OPTION_WEBRTC_CONFIGURATION = "webrtcConfiguration"
)

var (
	WithWebrtcConfiguration, GetWebrtcConfiguration = option.New[webrtc.Configuration](OPTION_WEBRTC_CONFIGURATION)
)
