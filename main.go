package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"kyleroberts.io/src/api"
	"kyleroberts.io/src/config"

	"github.com/joho/godotenv"
)

func apiLayer(appEnvironment api.AppEnvironment) {
	router := gin.Default()
	router.Use(appEnvironment.InboundRequestLog)
	routes := router.Group("/api")

	routes.GET("/ping", appEnvironment.Ping)

	routes.POST("/create", appEnvironment.CreateContainerGroup)
	routes.GET("/status", appEnvironment.ContainerGroupStatus)
	router.Run(fmt.Sprintf(":%d", appEnvironment.Config.API.Port))
}

func defineAppEnvironment() api.AppEnvironment {
	envLoadErr := godotenv.Load()
	if envLoadErr != nil {
		log.Panic("Failed to load env file data")
	}
	APP_CONF := config.GetAppConfig(false)
	LOGGER := initLogrus(APP_CONF.Logger)
	appEnvironment := api.AppEnvironment{
		Config: APP_CONF,
		Logger: LOGGER,
	}
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
