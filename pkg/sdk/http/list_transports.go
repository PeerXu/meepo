package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/sdk"
)

func (t *MeepoSDK) ListTransports() ([]*sdk.Transport, error) {
	var res http_api.ListTransportsResponse

	if err := t.doRequest("/v1/actions/list_transports", nil, &res, http.StatusOK); err != nil {
		return nil, err
	}

	return res.Transports, nil
}
