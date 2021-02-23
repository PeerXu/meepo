package http_sdk

import (
	"net/http"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
)

func (t *MeepoSDK) CloseTeleportation(name string) error {
	req := &http_api.CloseTeleportationRequest{
		Name: name,
	}

	if err := t.doRequest("/v1/actions/close_teleportation", req, nil, http.StatusNoContent); err != nil {
		return err
	}

	return nil
}
