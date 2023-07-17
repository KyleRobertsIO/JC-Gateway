package app

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"kyleroberts.io/src/api"
	"kyleroberts.io/src/config"
)

type Application struct {
	Name string
}

/*
Starts the entire application from
a single entrypoint.
*/
func (app *Application) Start() {
	// Define app environment to make accessible
	// configurations.
	appEnv, appEnvErr := app.defineEnvironment()
	if appEnvErr != nil {
		fmt.Println(appEnvErr.Error())
	}
	// Set the mode for the Gin web server
	if appEnv.Config.Gin.Mode == "release" {
		gin.SetMode("release")
	} else {
		gin.SetMode("debug")
	}
	appEnv.Logger.Info(
		fmt.Sprintf("Starting Gin Web Server in [%s] mode", gin.Mode()),
	)
	ginEngine := app.defineGinEngine(appEnv)
	ginEngine.Run(fmt.Sprintf(":%d", appEnv.Config.Gin.Port))
}

/*
Sets up the requirements for the Gin
web server
*/
func (app *Application) defineGinEngine(appEnv *api.AppEnvironment) *gin.Engine {
	// Define the Gin Engine
	engine := gin.Default()
	engine.Use(appEnv.MiddlewareInboundRequestLog)
	// Declare the Gin Engine routes
	routes := engine.Group("/api")
	routes.GET("/ping", appEnv.Ping)
	routes.POST("/create", appEnv.CreateContainerGroup)
	routes.GET("/status", appEnv.ContainerGroupStatus)
	return engine
}

/*
Sets up the requirements for the application
environment logger.
*/
func (app *Application) defineLogger(
	logConf config.LoggerConfig,
) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logLevel, logLevelErr := logrus.ParseLevel(
		string(logConf.LogLevel.LogLevelEnum),
	)
	if logLevelErr != nil {
		return nil, fmt.Errorf("failed establishing log level")
	}
	logger.SetLevel(logLevel)
	return logger, nil
}

/*
Sets up the requirements for the application
environment configuration.
*/
func (app *Application) defineEnvironment() (*api.AppEnvironment, error) {
	fmt.Println("loading application environment config")
	appConf := config.GetAppConfig(false)
	fmt.Println("initializing logger")
	logger, loggerErr := app.defineLogger(appConf.Logger)
	if loggerErr != nil {
		return nil, loggerErr
	}
	appEnv := api.AppEnvironment{
		Config: appConf,
		Logger: logger,
	}
	return &appEnv, nil
}
