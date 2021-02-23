package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	encoding_api "github.com/PeerXu/meepo/pkg/api/encoding"
)

type ListTransportsResponse struct {
	Transports []*encoding_api.Transport `json:"transports"`
}

func (s *HttpServer) ListTransports(c *gin.Context) {
	transports, err := s.meepo.ListTransports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	res := ListTransportsResponse{
		Transports: encoding_api.ConvertTransports(transports),
	}

	c.JSON(http.StatusOK, &res)
}
