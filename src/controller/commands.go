package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kyleroberts.io/src/azure"
)

type CreateContainerGroupInbound struct {
	ContainerGroupName string `json:"container_group_name"`
}

func CreateContainerGroup(context *gin.Context) {
	payload := new(CreateContainerGroupInbound)
	bindErr := context.BindJSON(&payload)
	if bindErr != nil {
		context.AbortWithError(http.StatusBadRequest, bindErr)
		return
	}
	fmt.Println(payload.ContainerGroupName)
	cg := azure.ContainerGroup{
		Subscription:  "abc",
		ResourceGroup: "temp-resource-group",
		Name:          "cron-container",
	}
	err := cg.Create("2022-09-01")
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{"message": "Created Azure Container Instance"},
		)
		return
	}
}
