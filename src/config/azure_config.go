package config

import (
	"fmt"
	"log"
	"os"
)

type AzureConfig struct {
	AuthType    string
	AuthDetails *AzureAuthDetails
}

type AzureAuthDetails struct {
	// ClientID is used for both service principal
	// and managed identity
	ClientID *string
	// Below are user for service principal
	ClientSecret *string
	TenantID     *string
}

func assembleServicePrincipalDetails() (*AzureAuthDetails, error) {
	clientID := os.Getenv("AZURE_AUTH_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_AUTH_CLIENT_SECRET")
	tenantID := os.Getenv("AZURE_AUTH_TENANT_ID")
	if clientID == "" {
		return nil, fmt.Errorf(
			"missing [AZURE_AUTH_CLIENT_ID] when using AZURE_AUTH_TYPE=SERVICE_PRINCIPAL",
		)
	}
	if clientSecret == "" {
		return nil, fmt.Errorf(
			"missing [AZURE_AUTH_CLIENT_SECRET] when using AZURE_AUTH_TYPE=SERVICE_PRINCIPAL",
		)
	}
	if tenantID == "" {
		return nil, fmt.Errorf(
			"missing [AZURE_AUTH_TENANT_ID] when using AZURE_AUTH_TYPE=SERVICE_PRINCIPAL",
		)
	}
	authDetails := AzureAuthDetails{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		TenantID:     &tenantID,
	}
	return &authDetails, nil
}

func assembleUserAssignedManagedIdentityDetails() (*AzureAuthDetails, error) {
	clientID := os.Getenv("AZURE_AUTH_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf(
			"missing [AZURE_AUTH_CLIENT_ID] when using AZURE_AUTH_TYPE=USER_ASSIGNED_MANAGED_IDENTITY",
		)
	}
	authDetails := AzureAuthDetails{
		ClientID:     &clientID,
		ClientSecret: nil,
		TenantID:     nil,
	}
	return &authDetails, nil
}

func determineAuthDetails(authType string) (*AzureAuthDetails, error) {
	switch authType {
	case "SERVICE_PRINCIPAL":
		return assembleServicePrincipalDetails()
	case "USER_ASSIGNED_MANAGED_IDENTITY":
		return assembleUserAssignedManagedIdentityDetails()
	default:
		return nil, fmt.Errorf("unknown [AZURE_AUTH_TYPE] value supplied")
	}
}

func assembleAzureConfig() AzureConfig {
	authType := os.Getenv("AZURE_AUTH_TYPE")
	if authType == "" {
		log.Fatal("missing required value for [AZURE_AUTH_TYPE]")
	}
	authDetails, authDetailsErr := determineAuthDetails(authType)
	if authDetailsErr != nil {
		log.Fatal(authDetailsErr.Error())
	}
	return AzureConfig{
		AuthType:    authType,
		AuthDetails: authDetails,
	}
}
