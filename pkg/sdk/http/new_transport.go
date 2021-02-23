package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/sdk"
)

func (t *MeepoSDK) NewTransport(peerID string) (*sdk.Transport, error) {
	var res http_api.NewTransportResponse
	req := &http_api.NewTransportRequest{
		PeerID: peerID,
	}

	if err := t.doRequest("/v1/actions/new_transport", req, &res, http.StatusCreated); err != nil {
		return nil, err
	}

	return res.Transport, nil
}
