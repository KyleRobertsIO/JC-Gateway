package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kyleroberts.io/src/api/payloads"
	"kyleroberts.io/src/azure/requests"
)

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponseError struct {
	Error ErrorDetails `json:"error"`
}

type ContainerGroupManager struct {
	AccessToken   string
	APIVersion    string
	Subscription  string
	ResourceGroup string
}

func (cgm *ContainerGroupManager) Create(payload *payloads.CreateContainerGroup) error {
	factory := ContainerGroupFactory{Payload: payload}
	reqBody := factory.Create()
	createErr := cgm.requestCreateGroup(reqBody)
	if createErr != nil {
		fmt.Println("failed to create container group")
	}
	return nil
}

func (cgm *ContainerGroupManager) encodeRequestBody(reqBody *requests.CreateContainerGroupBody) (*[]byte, error) {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to json encode request body")
	}
	return &b, nil
}

func (cgm *ContainerGroupManager) requestCreateGroup(
	reqBody requests.CreateContainerGroupBody,
) error {
	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroups/%s?api-version=%s",
		cgm.Subscription,
		cgm.ResourceGroup,
		reqBody.Name,
		cgm.APIVersion,
	)
	encodedReqBody, encodeErr := cgm.encodeRequestBody(&reqBody)
	if encodeErr != nil {
		return encodeErr
	}
	request, reqMakeErr := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(*encodedReqBody),
	)
	if reqMakeErr != nil {
		return fmt.Errorf("failed to create request body; %s", reqMakeErr.Error())
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cgm.AccessToken))
	request.Header.Set("Content-Type", "application/json")
	// Executing Request
	client := &http.Client{}
	response, reqErr := client.Do(request)
	if reqErr != nil {
		return fmt.Errorf("http request to %s failed; %s", request.URL, reqErr.Error())
	}
	defer response.Body.Close()
	// Decode Response Body
	bodyBytes, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return fmt.Errorf("failed to read response body; %s", bodyErr.Error())
	}
	// Check if successful result and return
	if response.StatusCode < 400 {
		return nil
	}
	// Create and error result and return
	var decodeTarget APIResponseError
	bodyDecodeErr := json.Unmarshal(bodyBytes, &decodeTarget)
	if bodyDecodeErr != nil {
		return fmt.Errorf("failed to decode error response body; %s", bodyDecodeErr.Error())
	}
	return nil
}
