package api

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/kylerobertsio/aci-job-manager/src/config"
	"github.com/sirupsen/logrus"
)

type AppEnvironment struct {
	Config          config.AppConfig
	Logger          *logrus.Logger
	AzureCredential *azidentity.DefaultAzureCredential
}
