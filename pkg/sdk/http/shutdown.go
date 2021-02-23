package http_sdk

import (
	"net/http"
)

func (t *MeepoSDK) Shutdown() error {
	return t.doRequest("/v1/actions/shutdown", nil, nil, http.StatusNoContent)
}
