package api

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"kyleroberts.io/src/api/payloads"
	"kyleroberts.io/src/azure"
)

/*
Authenticates with Azure Active Directory and returns a credential for Azure SDK libraries.
*/
func (env *AppEnvironment) AzureAuthenticate() {
	cred, credErr := azidentity.NewDefaultAzureCredential(nil)
	if credErr != nil {
		env.Logger.Info(
			fmt.Sprintf("failed to create azure credential for application; %s", credErr.Error()),
		)
	} else {
		env.AzureCredential = cred
	}
}

/*
Logs errors that take place inside of ContainerGroupManager command.
*/
func logContainerGroupManagerError(logger *logrus.Logger, err *azure.ContainerGroupManagerError) {
	logger.WithFields(logrus.Fields{
		"http_status_code": err.HttpStatusCode,
		"error_code":       err.Code,
		"error":            err.Error,
		"source":           err.Source,
	}).Warning("Outbound Response")
}

/*
Endpoint for creating or updating a Container Group
*/
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

/*
Endpoint for checking the status of a Container Group.
*/
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
