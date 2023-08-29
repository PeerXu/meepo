package rpc_simple_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *SimpleHttpServer) Routers() http.Handler {
	r := gin.New()
	v1Router := r.Group("/v1")
	v1ActionRouter := v1Router.Group("/actions")
	v1ActionRouter.POST("/simpleDo", s.HttpSimpleDo)
	v1ActionRouter.GET("/simpleDoStream", s.HttpSimpleDoStream)
	return r
}
