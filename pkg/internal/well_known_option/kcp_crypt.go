package well_known_option

import "github.com/PeerXu/meepo/pkg/internal/option"

const OPTION_KCP_CRYPT = "kcpCrypt"

var WithKcpCrypt, GetKcpCrypt = option.New[string](OPTION_KCP_CRYPT)
