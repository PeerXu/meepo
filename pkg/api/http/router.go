package http_api

import (
	"github.com/gin-gonic/gin"
)

func (s *HttpServer) getRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", s.Healthz)

	v1Action := router.Group("/v1/actions")
	v1Action.POST("/whoami", s.Whaomi)
	v1Action.POST("/shutdown", s.Shutdown)
	v1Action.POST("/teleport", s.Teleport)

	return router
}
