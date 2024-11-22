package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yairp7/go-common-lib/logger"
	"github.com/yairp7/papi"
	"github.com/yairp7/papi/config"
	"github.com/yairp7/papi/examples/simple-server/controllers"
)

type ConfigData struct {
	config.ServerConfig
	Names []string `json:"names"`
}

var loggerImpl logger.Logger = logger.NewMixedLogger(
	logger.WithLoggerLevel(logger.DEBUG),
	logger.WithLoggerImpl(logger.NewStdoutLogger(logger.DEBUG)),
	logger.WithLogSuffix("\n"),
)

func main() {
	godotenv.Load()
	config, err := config.ReadBase64Config[ConfigData](os.Getenv("CONFIG"))
	if err != nil {
		panic(err)
	}

	server := papi.NewServer(loggerImpl)
	server.Start(config.ServerConfig)
	namesController := controllers.NewNamesController(loggerImpl, config.Names)
	server.RegisterController(
		namesController,
		map[string][]papi.EndpointInfo{
			"/name": {
				{Method: http.MethodGet, Handler: namesController.Name},
				{Method: http.MethodPost, Handler: namesController.AddName},
			},
		},
	)
}
