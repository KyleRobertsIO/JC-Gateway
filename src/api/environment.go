package api

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"kyleroberts.io/src/config"
)

type AppEnvironment struct {
	Config          config.AppConfig
	AzureCredential *azidentity.ClientSecretCredential
}
