package papi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yairp7/go-common-lib/logger"
	"github.com/yairp7/papi/config"
	"github.com/yairp7/papi/controllers"
)

type EndpointInfo struct {
	Method      string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

type Server[ConfDataType any] struct {
	config     config.Config[ConfDataType]
	router     *Router
	loggerImpl logger.Logger
}

func NewServer[ConfDataType any](loggerImpl logger.Logger) *Server[ConfDataType] {
	return &Server[ConfDataType]{
		loggerImpl: loggerImpl,
	}
}

func (s *Server[ConfDataType]) Start(
	setup func(config config.Config[ConfDataType]),
) {
	godotenv.Load()
	config, err := config.ReadBase64Config[ConfDataType](os.Getenv("CONFIG"))
	if err != nil {
		panic(err)
	}

	if config.Env == "staging" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.config = config

	s.router, err = newRouter(s.loggerImpl)
	if err != nil {
		panic(err)
	}

	s.setupDefaultRoutes()
	setup(s.config)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "", config.Port),
		Handler: s.router.engine,
	}

	ctx, cancelSignalNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelSignalNotify()

	s.loggerImpl.Info("Server[v%s] ready on port %d\n", config.Version, config.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.loggerImpl.Error("Server error: %s\n", err)
		}
	}()

	<-ctx.Done()
	s.Shutdown(srv)
}

func (s *Server[ConfDataType]) RegisterController(
	controller controllers.Controller,
	endpoints map[string]EndpointInfo,
) {
	if controller == nil || endpoints == nil {
		return
	}

	for route, info := range endpoints {
		if info.Middlewares == nil || len(info.Middlewares) == 0 {
			s.router.engine.Handle(info.Method, route, info.Handler)
			continue
		}
		handlers := []gin.HandlerFunc{info.Handler}
		handlers = append(handlers, info.Middlewares...)
		s.router.engine.Handle(info.Method, route, handlers...)
	}
	s.router.registerController(controller)
}

func (s *Server[ConfDataType]) setupDefaultRoutes() {
	healthController := controllers.NewHealthController(s.loggerImpl)
	s.RegisterController(healthController, map[string]EndpointInfo{
		s.config.HeathCheckRoute: {Method: http.MethodGet, Handler: healthController.Status},
	})
}

func (s *Server[ConfDataType]) Shutdown(srv *http.Server) {
	s.loggerImpl.Info("Shutting server down...")

	s.router.shutdownRouter()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.loggerImpl.Error("Server shutdown failed: %s\n", err)
	}

	select {
	case <-ctx.Done():
		s.loggerImpl.Warn("Server shutdown timeout")
	default:
		s.loggerImpl.Info("Server shutdown complete")
	}
}
