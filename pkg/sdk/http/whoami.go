package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
)

func (t *MeepoSDK) Whoami() (string, error) {
	var res http_api.WhoamiResponse
	var err error

	if err = t.doRequest("/v1/actions/whoami", nil, &res, http.StatusOK); err != nil {
		return "", err
	}

	return res.ID, nil
}
