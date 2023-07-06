package config

import (
	"os"
)

type AzureConfig struct {
	TenantID string
}

func assembleAzureConfig() AzureConfig {
	return AzureConfig{
		TenantID: os.Getenv("AZURE.TENANT_ID"),
	}
}
