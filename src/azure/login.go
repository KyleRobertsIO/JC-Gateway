package azure

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AzureLoginBody struct {
	ClientID     string
	ClientSecret string
}

func AzLogin(
	tenantID string,
	body AzureLoginBody,
) error {
	// Define the HTTP client
	client := &http.Client{}
	// Define the request body
	data := url.Values{}
	data.Set("client_id", body.ClientID)
	data.Set("client_secret", body.ClientSecret)
	data.Set("scope", "https://management.azure.com")
	data.Set("grant_type", "./default")
	encodedData := data.Encode()
	// Define request URL
	url := fmt.Sprintf("login.microsoftonline.com/%s/oauth2/token", tenantID)
	// Create request struct
	request, err := http.NewRequest(http.MethodGet, url, strings.NewReader(encodedData))
	if err != nil {
		return fmt.Errorf("failed to assemble request struct")
	}
	// Append headers to request
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, reqErr := client.Do(request)
	if reqErr != nil {
		return fmt.Errorf("HTTP request to %s failed; %s", url, reqErr.Error())
	}
	return nil
}
