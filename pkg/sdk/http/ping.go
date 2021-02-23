package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
)

func (t *MeepoSDK) Ping(peerID string) error {
	req := &http_api.PingRequest{
		PeerID: peerID,
	}

	return t.doRequest("/v1/actions/ping", req, nil, http.StatusNoContent)
}
