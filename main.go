package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"kyleroberts.io/src/api"
	"kyleroberts.io/src/config"

	"github.com/joho/godotenv"
)

func apiLayer(appEnvironment api.AppEnvironment) {
	fmt.Println("Starting API Layer")
	router := gin.Default()
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
	appEnvironment := api.AppEnvironment{
		Config: APP_CONF,
	}
	return appEnvironment
}

func main() {
	appEnvironment := defineAppEnvironment()
	fmt.Println("Starting Application")
	apiLayer(appEnvironment)
}
