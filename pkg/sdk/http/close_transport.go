package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
)

func (t *MeepoSDK) CloseTransport(peerID string) error {
	req := &http_api.CloseTransportRequest{
		PeerID: peerID,
	}

	return t.doRequest("/v1/actions/close_transport", req, nil, http.StatusNoContent)
}
