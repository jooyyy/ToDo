package server

import "github.com/gin-gonic/gin"

func (s *Service) initRouter(r gin.IRouter) {
	r.GET("/api/home/list", s.GetHomeList)
}