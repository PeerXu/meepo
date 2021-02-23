package http_api

import (
	"github.com/gin-gonic/gin"
)

func (s *HttpServer) getRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", s.Healthz)

	v1Action := router.Group("/v1/actions")
	v1Action.POST("/version", s.Version)
	v1Action.POST("/ping", s.Ping)
	v1Action.POST("/whoami", s.Whaomi)
	v1Action.POST("/shutdown", s.Shutdown)
	v1Action.POST("/teleport", s.Teleport)

	v1Action.POST("/new_transport", s.NewTransport)
	v1Action.POST("/close_transport", s.CloseTransport)
	v1Action.POST("/list_transports", s.ListTransports)

	v1Action.POST("/new_teleportation", s.NewTeleportation)
	v1Action.POST("/close_teleportation", s.CloseTeleportation)
	v1Action.POST("/list_teleportations", s.ListTeleportations)

	return router
}
