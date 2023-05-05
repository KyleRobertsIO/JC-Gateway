package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"kyleroberts.io/src/config"
	"kyleroberts.io/src/controller"
)

func apiLayer(port int) {
	fmt.Println("Starting API Layer")
	router := gin.Default()
	routes := router.Group("/api")
	routes.GET("/ping", controller.Ping)

	routes.POST("/create", controller.CreateContainerGroup)
	router.Run(fmt.Sprintf(":%d", port))
}

func main() {
	APP_CONF := config.GetAppConfig(false)
	fmt.Println("Starting Application")
	apiLayer(APP_CONF.API.Port)
}
