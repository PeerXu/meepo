package http_api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *HttpServer) Shutdown(c *gin.Context) {
	go func() {
		time.Sleep(1 * time.Second)
		s.Stop(context.TODO())
	}()
	c.Writer.WriteHeader(http.StatusNoContent)
}
