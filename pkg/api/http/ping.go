package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingRequest struct {
	PeerID  string `json:"peerID" binding:"required"`
	Payload string `json:"payload,omitempty"`
}

func (s *HttpServer) Ping(c *gin.Context) {
	var err error
	var req PingRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	if err = s.meepo.Ping(req.PeerID, req.Payload); err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
