package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type AzureAuthTypeEnum string

const (
	SERVICE_PRINCIPAL                AzureAuthTypeEnum = "SERVICE_PRINCIPAL"
	USER_ASSIGNED_MANAGED_IDENTITY   AzureAuthTypeEnum = "USER_ASSIGNED_MANAGED_IDENTITY"
	SYSTEM_ASSIGNED_MANAGED_IDENTITY AzureAuthTypeEnum = "SYSTEM_ASSIGNED_MANAGED_IDENTITY"
)

type AzureAuthType struct {
	AzureAuthTypeEnum
}

func (aat *AzureAuthType) FromStr(typeStr string) error {
	switch strings.ToUpper(typeStr) {
	case "SERVICE_PRINCIPAL":
		aat.AzureAuthTypeEnum = SERVICE_PRINCIPAL
	case "USER_ASSIGNED_MANAGED_IDENTITY":
		aat.AzureAuthTypeEnum = USER_ASSIGNED_MANAGED_IDENTITY
	case "SYSTEM_ASSIGNED_MANAGED_IDENTITY":
		aat.AzureAuthTypeEnum = SYSTEM_ASSIGNED_MANAGED_IDENTITY
	default:
		allowedOptions := "['SERVICE_PRINCIPAL', 'USER_ASSIGNED_MANAGED_IDENTITY', 'SYSTEM_ASSIGNED_MANAGED_IDENTITY']"
		msg := "failed to determine provided AuthType from string;"
		errMsg := fmt.Sprintf("%s acceptable values are %s", msg, allowedOptions)
		return errors.New(errMsg)
	}
	return nil
}

type AzureAppConfig struct {
	ClientID     string
	ClientSecret string
	Scope        string
	AuthType     AzureAuthType
}

type AzureConfig struct {
	TenantID string
	App      AzureAppConfig
}

func assembleAzureAppConfig() AzureAppConfig {
	authType := AzureAuthType{}
	authTypeErr := authType.FromStr(os.Getenv("AZURE.APP.AUTH_TYPE"))
	if authTypeErr != nil {
		log.Fatal(authTypeErr.Error())
	}
	return AzureAppConfig{
		ClientID:     os.Getenv("AZURE.APP.CLIENT_ID"),
		ClientSecret: os.Getenv("AZURE.APP.CLIENT_SECRET"),
		Scope:        os.Getenv("AZURE.APP.SCOPE"),
		AuthType:     authType,
	}
}

func assembleAzureConfig() AzureConfig {
	return AzureConfig{
		TenantID: os.Getenv("AZURE.TENANT_ID"),
		App:      assembleAzureAppConfig(),
	}
}
