package api

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/sirupsen/logrus"
	"kyleroberts.io/src/config"
)

type AppEnvironment struct {
	Config          config.AppConfig
	Logger          *logrus.Logger
	AzureCredential *azidentity.ClientSecretCredential
}
