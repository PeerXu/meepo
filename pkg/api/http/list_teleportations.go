package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	encoding_api "github.com/PeerXu/meepo/pkg/api/encoding"
)

type ListTeleportationsResponse struct {
	Teleportations []*encoding_api.Teleportation `json:"teleportations"`
}

func (s *HttpServer) ListTeleportations(c *gin.Context) {
	teleportations, err := s.meepo.ListTeleportations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	res := ListTeleportationsResponse{
		Teleportations: encoding_api.ConvertTeleportations(teleportations),
	}

	c.JSON(http.StatusOK, &res)
}
