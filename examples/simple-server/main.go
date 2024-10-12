package main

import (
	"net/http"

	"github.com/yairp7/go-common-lib/logger"
	"github.com/yairp7/papi"
	"github.com/yairp7/papi/config"
	"github.com/yairp7/papi/examples/simple-server/controllers"
)

type ConfigData struct {
	Names []string `json:"names"`
}

var loggerImpl logger.Logger = logger.NewMixedLogger(
	logger.WithLoggerLevel(logger.DEBUG),
	logger.WithLoggerImpl(logger.NewStdoutLogger(logger.DEBUG)),
	logger.WithLogSuffix("\n"),
)

func main() {
	server := papi.NewServer[ConfigData](loggerImpl)
	server.Start(func(conf config.Config[ConfigData]) {
		namesController := controllers.NewNamesController(loggerImpl, conf.Data.Names)
		server.RegisterController(
			namesController,
			map[string][]papi.EndpointInfo{
				"/name": {
					{Method: http.MethodGet, Handler: namesController.Name},
					{Method: http.MethodPost, Handler: namesController.AddName},
				},
			},
		)
	})
}
