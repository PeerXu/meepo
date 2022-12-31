package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_WEBRTC_RECEIVE_BUFFER_SIZE = "webrtcReceiveBufferSize"

var WithWebrtcReceiveBufferSize, GetWebrtcReceiveBufferSize = option.New[uint32](OPTION_WEBRTC_RECEIVE_BUFFER_SIZE)
