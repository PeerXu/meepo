package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/sdk"
)

func (t *MeepoSDK) ListTeleportations() ([]*sdk.Teleportation, error) {
	var res http_api.ListTeleportationsResponse

	if err := t.doRequest("/v1/actions/list_teleportations", nil, &res, http.StatusOK); err != nil {
		return nil, err
	}

	return res.Teleportations, nil
}
