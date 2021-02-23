package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CloseTeleportationRequest struct {
	Name string `json:"name" binding:"required"`
}

func (s *HttpServer) CloseTeleportation(c *gin.Context) {
	var err error
	var req CloseTeleportationRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ParseError(err))
		return
	}

	if err = s.meepo.CloseTeleportation(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, ParseError(err))
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
