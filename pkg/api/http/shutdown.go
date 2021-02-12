package http_api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *HttpServer) Shutdown(c *gin.Context) {
	go func() {
		s.Stop(context.TODO())
	}()
	c.Writer.WriteHeader(http.StatusNoContent)
}
