package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/sdk"
)

func (t *MeepoSDK) Version() (*sdk.Version, error) {
	var res http_api.VersionResponse

	if err := t.doRequest("/v1/actions/version", nil, &res, http.StatusOK); err != nil {
		return nil, err
	}

	return &res, nil
}
