package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type AzureTokenAuthRequirements struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	Scope        string
}

func determineAuthType() {

}

func GetAzureToken(
	requirements AzureTokenAuthRequirements,
) (*azidentity.ClientSecretCredential, error) {
	cred, credErr := azidentity.NewClientSecretCredential(
		requirements.TenantID,
		requirements.ClientID,
		requirements.ClientSecret,
		nil,
	)
	if credErr != nil {
		return nil, fmt.Errorf("failed to collect Azure access credential; %s", credErr.Error())
	}
	return cred, nil
}
