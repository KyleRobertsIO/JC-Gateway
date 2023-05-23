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

type ErrorSource string

const (
	AzureResourceManager ErrorSource = "AZURE_RESOURCE_MANAGER"
	Proxy                ErrorSource = "PROXY"
)

type ContainerGroupManagerError struct {
	Error          string      `json:"error"`
	Source         ErrorSource `json:"source"`
	Code           string      `json:"error_code"`
	HttpStatusCode int         `json:"-"`
}

type ContainerGroupManager struct {
	AccessToken   string
	APIVersion    string
	Subscription  string
	ResourceGroup string
}

func (cgm *ContainerGroupManager) Create(
	payload *payloads.CreateContainerGroup,
) *ContainerGroupManagerError {
	factory := ContainerGroupFactory{Payload: payload}
	reqBody := factory.Create()
	return cgm.requestCreateGroup(reqBody)
}

func (cgm *ContainerGroupManager) encodeRequestBody(
	reqBody *requests.CreateContainerGroupBody,
) (*[]byte, *ContainerGroupManagerError) {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to encode request body; %s", err.Error()),
			Source:         Proxy,
			Code:           "REQUEST_BODY_ENCODE_ERROR",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}
	return &b, nil
}

func (cgm *ContainerGroupManager) requestCreateGroup(
	reqBody requests.CreateContainerGroupBody,
) *ContainerGroupManagerError {
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
		return &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to build request for Azure REST API; %s", reqMakeErr.Error()),
			Source:         Proxy,
			Code:           "REQUEST_BUILD_ERROR",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cgm.AccessToken))
	request.Header.Set("Content-Type", "application/json")
	// Executing Request
	client := &http.Client{}
	response, reqErr := client.Do(request)
	if reqErr != nil {
		return &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to issue request to Azure REST API; %s", reqErr.Error()),
			Source:         Proxy,
			Code:           "CREATE_REQUEST_ERROR",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}
	defer response.Body.Close()
	// Decode Response Body
	bodyBytes, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to read response body of request; %s", bodyErr.Error()),
			Source:         Proxy,
			Code:           "RESPONSE_BODY_READ_ERROR",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}
	// Check if successful result and return
	if response.StatusCode < 400 {
		return nil
	}
	// Create and error result and return
	var decodeTarget APIResponseError
	bodyDecodeErr := json.Unmarshal(bodyBytes, &decodeTarget)
	if bodyDecodeErr != nil {
		return &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to decode response body of bad request; %s", bodyDecodeErr.Error()),
			Source:         Proxy,
			Code:           "BAD_REQUEST_BODY_DECODE_ERROR",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}
	return &ContainerGroupManagerError{
		Error:          fmt.Sprintf(decodeTarget.Error.Message),
		Source:         AzureResourceManager,
		Code:           decodeTarget.Error.Code,
		HttpStatusCode: http.StatusBadRequest,
	}
}
