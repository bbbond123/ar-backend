package router

import (
	"github.com/gin-gonic/gin"
)

// RouteRegister 路由注册器接口
type RouteRegister interface {
	Register(r *gin.RouterGroup)
}

var routeRegisters []RouteRegister

// Register 注册路由模块
func Register(rr RouteRegister) {
	routeRegisters = append(routeRegisters, rr)
}

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	// 注册所有模块路由
	for _, rr := range routeRegisters {
		rr.Register(api)
	}

	return r
}
