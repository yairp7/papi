package papi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yairp7/go-common-lib/logger"
	"github.com/yairp7/papi/controllers"
)

type Router struct {
	controllers []controllers.Controller
	isActive    bool
	loggerImpl  logger.Logger
	engine      *gin.Engine
}

func (r *Router) isServerActive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !r.isActive {
			ctx.AbortWithError(http.StatusServiceUnavailable, serviceUnavailableError)
			return
		}
	}
}

func newRouter(loggerImpl logger.Logger) (*Router, error) {
	router := Router{
		loggerImpl:  loggerImpl,
		controllers: make([]controllers.Controller, 0),
		isActive:    true,
	}

	router.engine = gin.New()
	router.engine.Use(gin.Logger())
	router.engine.Use(gin.Recovery())
	router.engine.Use(router.isServerActive())

	return &router, nil
}

func (r *Router) registerController(controller controllers.Controller) {
	r.controllers = append(r.controllers, controller)
}

func (r *Router) shutdownRouter() {
	r.loggerImpl.Info("Router Shutdown")
	r.isActive = false
	for _, closeable := range r.controllers {
		closeable.Close()
	}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
