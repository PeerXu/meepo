package sdk

import (
	"fmt"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
)

var ExtractError = http_api.ExtractError
var UnimplementedError = fmt.Errorf("Unimplemented")

func UnsupportedMeepoSDKDriverError(name string) error {
	return fmt.Errorf("Unsupported MeepoSDK driver: %s", name)
}
