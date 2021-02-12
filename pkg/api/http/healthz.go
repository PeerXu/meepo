package http_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*HttpServer) Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
