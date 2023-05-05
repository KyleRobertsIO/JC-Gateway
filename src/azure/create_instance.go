package azure

import (
	"fmt"
	"net/http"
)

type ContainerGroup struct {
	Subscription  string
	ResourceGroup string
	Name          string
}

func (cg *ContainerGroup) Create(apiVersion string) error {
	client := &http.Client{}
	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroups/%s?api-version=%s",
		cg.Subscription,
		cg.ResourceGroup,
		cg.Name,
		apiVersion,
	)
	request, err := http.NewRequest(
		http.MethodPut,
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to create container group.")
	}
	_, reqErr := client.Do(request)
	if reqErr != nil {
		return fmt.Errorf("HTTP request to %s failed; %s", url, reqErr.Error())
	}
	return nil
}
