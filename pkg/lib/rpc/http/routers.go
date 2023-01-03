package rpc_http

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *HttpServer) Routers() http.Handler {
	r := gin.New()
	r.Use(cors.Default())
	v1Router := r.Group("/v1")
	v1ActionRouter := v1Router.Group("/actions")
	v1ActionRouter.POST("/do", s.HttpDo)
	return r
}
