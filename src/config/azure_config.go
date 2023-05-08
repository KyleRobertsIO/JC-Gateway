package config

import "os"

type AzureAppConfig struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

type AzureConfig struct {
	TenantID string
	App      AzureAppConfig
}

func assembleAzureAppConfig() AzureAppConfig {
	return AzureAppConfig{
		ClientID:     os.Getenv("AZURE.APP.CLIENT_ID"),
		ClientSecret: os.Getenv("AZURE.APP.CLIENT_SECRET"),
		Scope:        os.Getenv("AZURE.APP.SCOPE"),
	}
}

func assembleAzureConfig() AzureConfig {
	return AzureConfig{
		TenantID: os.Getenv("AZURE.TENANT_ID"),
		App:      assembleAzureAppConfig(),
	}
}
