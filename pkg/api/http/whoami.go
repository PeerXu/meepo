package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WhoamiResponse struct {
	ID string `json:"id"`
}

func (s *HttpServer) Whaomi(c *gin.Context) {
	res := &WhoamiResponse{
		ID: s.meepo.ID(),
	}

	c.JSON(http.StatusOK, res)
}
