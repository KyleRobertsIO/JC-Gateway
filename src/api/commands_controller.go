package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	authRes, loginErr := azure.GetAzureToken(authRequirements)
	if loginErr != nil {
		fmt.Println("failed to authenticate with Azure")
	}
	env.AzureAccessToken = authRes.AccessToken
}

type SubnetDetails struct {
	VNetName      string `json:"vnet_name"`
	SubnetName    string `json:"subnet_name"`
	Subscription  string `json:"subscription"`
	ResourceGroup string `json:"resource_group"`
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
		AccessToken:   env.AzureAccessToken,
		APIVersion:    "2022-09-01",
		Subscription:  payload.Subscription,
		ResourceGroup: payload.ResourceGroup,
	}
	createErr := cgManager.Create(payload)
	if createErr != nil {
		context.JSON(
			createErr.HttpStatusCode,
			createErr,
		)
		return
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{"message": "Created Azure Container Instance"},
		)
		return
	}
}
