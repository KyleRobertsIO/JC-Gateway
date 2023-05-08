package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kyleroberts.io/src/azure"
	"kyleroberts.io/src/templates"
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

type CreateContainerGroupInbound struct {
	Subscription       string `json:"subscription"`
	ResourceGroup      string `json:"resource_group"`
	ContainerGroupName string `json:"container_group_name"`
	TemplateName       string `json:"template_name"`
}

func (env *AppEnvironment) CreateContainerGroup(context *gin.Context) {
	// env.AzureAuthenticate()
	payload := new(CreateContainerGroupInbound)
	bindErr := context.BindJSON(&payload)
	if bindErr != nil {
		context.AbortWithError(http.StatusBadRequest, bindErr)
		return
	}
	templateConfig, templateErr := templates.Parse(payload.TemplateName)
	if templateErr != nil {
		context.AbortWithError(http.StatusBadRequest, templateErr)
		return
	}
	fmt.Println(templateConfig.Image)
	cg := azure.ContainerGroup{
		Subscription:  payload.Subscription,
		ResourceGroup: payload.ResourceGroup,
		Name:          payload.ContainerGroupName,
	}
	err := cg.Create("2022-09-01", env.AzureAccessToken)
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
