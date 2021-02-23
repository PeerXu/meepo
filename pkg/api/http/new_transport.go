package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	encoding_api "github.com/PeerXu/meepo/pkg/api/encoding"
	"github.com/PeerXu/meepo/pkg/transport"
)

type NewTransportRequest struct {
	PeerID string `json:"peerID" binding:"required"`
}

type NewTransportResponse struct {
	Transport *encoding_api.Transport `json:"transport"`
}

func (s *HttpServer) NewTransport(c *gin.Context) {
	var err error
	var req NewTransportRequest
	var tp transport.Transport

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	if tp, err = s.meepo.NewTransport(req.PeerID); err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	c.JSON(http.StatusCreated, encoding_api.ConvertTransport(tp))
}
