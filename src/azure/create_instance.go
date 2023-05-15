package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"kyleroberts.io/src/templates"
)

//##############################
// Outbound Azure Payload
//##############################

type ResourceRequest struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memoryInGB"`
	// Will need to work in GPU allocation at some point
}

type ContainerResources struct {
	ResourceRequest ResourceRequest `json:"requests"`
}

type ContainerProperties struct {
	Command              []string                 `json:"command"`
	EnvironmentVariables []string                 `json:"environmentVariables"`
	Image                string                   `json:"image"`
	Ports                []map[string]interface{} `json:"ports"`
	Resources            ContainerResources       `json:"resources"`
}

type Container struct {
	Name       string              `json:"name"`
	Properties ContainerProperties `json:"properties"`
}

type RequestProperties struct {
	Containers []Container `json:"container"`
}

type CreateRequestBody struct {
	Location   string            `json:"location"`
	Properties RequestProperties `json:"properties"`
}

//##############################
// Inbound Client Payload
//##############################

type ContainerGroup struct {
	Subscription  string
	ResourceGroup string
	Name          string
}

//##############################
// API Responses
//##############################

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponseError struct {
	Error ErrorDetails `json:"error"`
}

/*
Converts the provided port definition from a yaml template to a seperate variables.

Example:

"80/tcp" --> 80, "tcp"
*/
func (cg *ContainerGroup) explodePort(portDefinition string) (int, string, error) {
	portDetails := strings.Split(portDefinition, "/")
	portNumber, strConverErr := strconv.Atoi(string(portDetails[0]))
	if strConverErr != nil {
		return 0, "", fmt.Errorf(
			"failed to parse port number from string; %s",
			strConverErr.Error(),
		)
	}
	portProtocol := portDetails[1]
	return portNumber, portProtocol, nil
}

/*
Converts a provided collection of port definition strings to a seperated map result.

Example:

"80/tcp" --> { "port": 80, "protocol": "tcp" }
*/
func (cg *ContainerGroup) translatePortsToMaps(ports []string) ([]map[string]interface{}, error) {
	var translatedPorts []map[string]interface{}
	for _, port := range ports {
		tempMap := make(map[string]interface{})
		portNum, protocol, expodeErr := cg.explodePort(port)
		if expodeErr != nil {
			return nil, fmt.Errorf("container port translation issue; %s", expodeErr.Error())
		}
		tempMap["port"] = portNum
		tempMap["protocol"] = protocol
		translatedPorts = append(translatedPorts, tempMap)
	}
	return translatedPorts, nil
}

func (cg *ContainerGroup) buildCreateRequestBody(
	templateConfig templates.YamlConfig,
) (*CreateRequestBody, error) {
	portMaps, portTranslateErr := cg.translatePortsToMaps(templateConfig.Ports)
	if portTranslateErr != nil {
		return nil, fmt.Errorf(
			"container group build create request failed; %s",
			portTranslateErr.Error(),
		)
	}
	containerResources := ContainerResources{
		ResourceRequest: ResourceRequest{
			CPU:    templateConfig.Resources.CPU,
			Memory: templateConfig.Resources.Memory,
		},
	}
	container := Container{
		Name: "sample-container-name-1",
		Properties: ContainerProperties{
			Command:              templateConfig.Commands,
			EnvironmentVariables: templateConfig.EnvironmentVariables,
			Image:                templateConfig.Image,
			Ports:                portMaps,
			Resources:            containerResources,
		},
	}
	var containerArr = []Container{}
	containerArr = append(containerArr, container)
	body := CreateRequestBody{
		Location: "canada-central",
		Properties: RequestProperties{
			Containers: containerArr,
		},
	}
	return &body, nil
}

func (cg *ContainerGroup) azExecuteCreateRequest(request *http.Request) (*http.Response, error) {
	// Send HTTP request
	client := &http.Client{}
	response, reqErr := client.Do(request)
	if reqErr != nil {
		return nil, fmt.Errorf("http request to %s failed; %s", request.URL, reqErr.Error())
	}
	defer response.Body.Close()
	bodyBytes, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return nil, fmt.Errorf("failed to read response body; %s", bodyErr.Error())
	}
	// Check if successful result and return
	if response.StatusCode < 400 {
		return response, nil
	}
	// Create and error result and return
	var decodeTarget APIResponseError
	bodyDecodeErr := json.Unmarshal(bodyBytes, &decodeTarget)
	if bodyDecodeErr != nil {
		return nil, fmt.Errorf("failed to decode error response body; %s", bodyDecodeErr.Error())
	}
	return nil, fmt.Errorf("failed to create container group; %s", decodeTarget.Error.Message)
}

func (cg *ContainerGroup) encodeCreateBody(reqBody *CreateRequestBody) (*[]byte, error) {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to json encode request body")
	}
	return &b, nil
}

func (cg *ContainerGroup) Create(
	apiVersion string,
	accessToken string,
	templateConfig templates.YamlConfig,
) error {
	url := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroups/%s?api-version=%s",
		cg.Subscription,
		cg.ResourceGroup,
		cg.Name,
		apiVersion,
	)
	reqBody, reqBodyErr := cg.buildCreateRequestBody(templateConfig)
	if reqBodyErr != nil {
		return reqBodyErr
	}
	fmt.Println(reqBody)
	encodedReqBody, encodeErr := cg.encodeCreateBody(reqBody)
	if encodeErr != nil {
		return encodeErr
	}
	request, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(*encodedReqBody),
	)
	if err != nil {
		return fmt.Errorf("Failed to create container group.")
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")
	res, resErr := cg.azExecuteCreateRequest(request)
	if resErr != nil {
		fmt.Println(resErr.Error())
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bodyBytes))
	return nil
}
