package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"kyleroberts.io/src/api"
	"kyleroberts.io/src/config"
)

func apiLayer(appEnvironment api.AppEnvironment) {
	router := gin.Default()
	router.Use(appEnvironment.InboundRequestLog)
	routes := router.Group("/api")

	routes.GET("/ping", appEnvironment.Ping)

	routes.POST("/create", appEnvironment.CreateContainerGroup)
	routes.GET("/status", appEnvironment.ContainerGroupStatus)
	router.Run(fmt.Sprintf(":%d", appEnvironment.Config.Gin.Port))
}

func defineAppEnvironment() api.AppEnvironment {
	fmt.Println("loading application config")
	APP_CONF := config.GetAppConfig(false)
	fmt.Println("initializing logger")
	LOGGER := initLogrus(APP_CONF.Logger)
	appEnvironment := api.AppEnvironment{
		Config: APP_CONF,
		Logger: LOGGER,
	}
	fmt.Println("application environment created")
	return appEnvironment
}

func initLogrus(logConf config.LoggerConfig) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logLevel, logLevelErr := logrus.ParseLevel(string(logConf.LogLevel.LogLevelEnum))
	if logLevelErr != nil {
		fmt.Println("failed establish log level")
		return nil
	}
	logger.SetLevel(logLevel)
	return logger
}

func main() {
	appEnvironment := defineAppEnvironment()
	if appEnvironment.Config.Gin.Mode == "release" {
		gin.SetMode("release")
	} else {
		gin.SetMode("debug")
	}
	appEnvironment.Logger.Info(
		fmt.Sprintf("Starting Gin Web Server in [%s] mode", gin.Mode()),
	)
	apiLayer(appEnvironment)
}