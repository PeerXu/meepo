package well_known_option

import (
	"net/http"

	"github.com/PeerXu/meepo/pkg/internal/option"
)

const (
	OPTION_HTTP_CLIENT = "httpClient"
)

var (
	WithHttpClient, GetHttpClient = option.New[*http.Client](OPTION_HTTP_CLIENT)
)
