package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AzureTokenAuthRequirements struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	Scope        string
}

type AzureLoginBody struct {
	ClientID     string
	ClientSecret string
}

type AzureTokenLoginResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

func createLoginRequest(requirements AzureTokenAuthRequirements) (*http.Request, error) {
	// Define the request body
	data := url.Values{}
	data.Set("client_id", requirements.ClientID)
	data.Set("client_secret", requirements.ClientSecret)
	// data.Set("scope", fmt.Sprintf("%s/.default", requirements.Scope))
	data.Set("scope", "https://management.azure.com/.default")
	data.Set("grant_type", "client_credentials")
	encodedData := data.Encode()
	// Define request URL
	url := fmt.Sprintf(
		"https://login.microsoftonline.com/%s/oauth2/v2.0/token",
		requirements.TenantID,
	)
	// Create request struct
	request, reqErr := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(encodedData),
	)
	if reqErr != nil {
		return nil, fmt.Errorf("failed to assemble request struct")
	}
	// Append headers to request
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}

func GetAzureToken(requirements AzureTokenAuthRequirements) (*AzureTokenLoginResponse, error) {
	// Define the HTTP client
	client := &http.Client{}
	// Define the HTTP request
	request, reqInitErr := createLoginRequest(requirements)
	if reqInitErr != nil {
		return nil, reqInitErr
	}
	res, reqErr := client.Do(request)
	if reqErr != nil {
		return nil, fmt.Errorf(
			"HTTP request to %s failed; %s",
			request.URL,
			reqErr.Error(),
		)
	}
	defer res.Body.Close()
	bodyBytes, bodyErr := ioutil.ReadAll(res.Body)
	if bodyErr != nil {
		fmt.Println("Response body failed to be parsed")
	}
	result := new(AzureTokenLoginResponse)
	jsonParseErr := json.Unmarshal(bodyBytes, result)
	if jsonParseErr != nil {
		return nil, jsonParseErr
	}
	return result, nil
}
