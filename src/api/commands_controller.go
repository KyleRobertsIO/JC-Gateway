package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"kyleroberts.io/src/api/payloads"
	"kyleroberts.io/src/azure"
)

func (env *AppEnvironment) AzureAuthenticate() {
	authRequirements := azure.AzureTokenAuthRequirements{
		ClientID:     env.Config.Azure.App.ClientID,
		ClientSecret: env.Config.Azure.App.ClientSecret,
		TenantID:     env.Config.Azure.TenantID,
		Scope:        env.Config.Azure.App.Scope,
	}
	cred, credErr := azure.GetAzureToken(authRequirements)
	if credErr != nil {
		fmt.Println(credErr.Error())
	} else {
		env.AzureCredential = cred
	}
}

type SubnetDetails struct {
	VNetName      string `json:"vnet_name"`
	SubnetName    string `json:"subnet_name"`
	Subscription  string `json:"subscription"`
	ResourceGroup string `json:"resource_group"`
}

func logContainerGroupManagerError(logger *logrus.Logger, err *azure.ContainerGroupManagerError) {
	logger.WithFields(logrus.Fields{
		"http_status_code": err.HttpStatusCode,
		"error_code":       err.Code,
		"error":            err.Error,
		"source":           err.Source,
	}).Warning("Outbound Response")
}

func (env *AppEnvironment) CreateContainerGroup(context *gin.Context) {
	env.AzureAuthenticate()
	payload := new(payloads.CreateContainerGroup)
	bindErr := context.BindJSON(&payload)
	if bindErr != nil {
		context.AbortWithError(http.StatusBadRequest, bindErr)
		return
	}
	cgManager := azure.ContainerGroupManager{
		Credential:    env.AzureCredential,
		Subscription:  payload.Subscription,
		ResourceGroup: payload.ResourceGroup,
	}
	createErr := cgManager.CreateOrUpdate(payload)
	if createErr != nil {
		logContainerGroupManagerError(env.Logger, createErr)
		context.JSON(createErr.HttpStatusCode, createErr)
		return
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{"message": "Issued For Azure Container Instance"},
		)
		return
	}
}

func (env *AppEnvironment) ContainerGroupStatus(context *gin.Context) {
	env.AzureAuthenticate()
	cgManager := azure.ContainerGroupManager{
		Credential:    env.AzureCredential,
		Subscription:  context.Query("subscription"),
		ResourceGroup: context.Query("resource_group"),
	}
	containerStatus, statusErr := cgManager.Status(context.Query("group_name"))
	if statusErr != nil {
		logContainerGroupManagerError(env.Logger, statusErr)
		context.JSON(statusErr.HttpStatusCode, statusErr)
		return
	} else {
		context.JSON(http.StatusOK, containerStatus)
		return
	}
}
