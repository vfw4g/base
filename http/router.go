package vfwhttp

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

type (
	HttpRouter []Route
)

var (
	httpRouters = make([]HttpRouter, 0, 0)

	middlewares = []gin.HandlerFunc{
		gin.Logger(),
		//middleware.Recovery(),
		//middleware.ValidateTokenHandlerFunc(),
		CorsMW(),
	}
)

func CorsMW() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{"http://localhost:8082"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS", "DELETE", "HEAD"},
		AllowHeaders:     []string{"authorization", "content-type"},
		AllowCredentials: true,
		// MaxAge: 12 * time.Hour,
	})
}

type GinEngin struct {
	*gin.Engine
}

func DefaultGin() *GinEngin {
	return &GinEngin{
		gin.Default(),
	}
}

func New() *GinEngin {
	return &GinEngin{
		gin.New(),
	}
}

func (e *GinEngin) NewGroup(s string) *GinRouterGroup {
	g := e.Group(s)
	return &GinRouterGroup{g}
}

type GinRouterGroup struct {
	*gin.RouterGroup
}

func (grg *GinRouterGroup) RegisterHttpRouters(hrs ...HttpRouter) {
	for _, hr := range hrs {
		for _, r := range hr {
			for _, v := range r.Methods {
				grg.Handle(v, r.Pattern, r.HandlerFunc)
			}
		}
	}
}

func (grg *GinRouterGroup) RegisterMiddleWare(r *gin.RouterGroup) {
	for _, middleware := range middlewares {
		grg.Use(middleware)
	}
}
