package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CloseTransportRequest struct {
	PeerID string `json:"peerID" binding:"required"`
}

func (s *HttpServer) CloseTransport(c *gin.Context) {
	var err error
	var req CloseTransportRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	if err = s.meepo.CloseTransport(req.PeerID); err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
