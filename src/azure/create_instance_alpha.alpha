package azure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"kyleroberts.io/src/azure/requests"
	"kyleroberts.io/src/templates"
)

//##############################
//
// Inbound Client Payload
//
//##############################

type ContainerGroupSubnet struct {
	VNetName      string
	SubnetName    string
	ResourceGroup string
	Subscription  string
}

type ContainerGroup struct {
	Subscription  string
	ResourceGroup string
	Name          string
	Subnet        ContainerGroupSubnet
}

//##############################
//
// API Responses
//
//##############################

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponseError struct {
	Error ErrorDetails `json:"error"`
}

func (cg *ContainerGroup) buildSubnetId(
	subscriptionId string,
	resourceGroup string,
	vnetName string,
	subnetName string,
) string {
	subnetId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s",
		subscriptionId,
		resourceGroup,
		vnetName,
		subnetName,
	)
	return subnetId
}

func (cg *ContainerGroup) translateOSType(os string) (string, error) {
	switch strings.ToUpper(os) {
	case "LINUX":
		return "Linux", nil
	case "WINDOWS":
		return "Windows", nil
	default:
		allowedOptions := "['LINUX', 'WINDOWS']"
		msg := "failed to determine provided OSType from string;"
		errMsg := fmt.Sprintf("%s acceptable values are %s", msg, allowedOptions)
		return "", errors.New(errMsg)
	}
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

func (cg *ContainerGroup) translateEnvironmentVars(
	envVars []templates.YamlEnvironmentVariable,
) []map[string]interface{} {
	var translatedVars []map[string]interface{}
	for _, yamlVar := range envVars {
		tempMap := make(map[string]interface{})
		tempMap["name"] = yamlVar.Name
		if yamlVar.Secure {
			tempMap["secureValue"] = yamlVar.Value
		} else {
			tempMap["value"] = yamlVar.Value
		}
		translatedVars = append(translatedVars, tempMap)
	}
	return translatedVars
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
) (*requests.CreateContainerGroupBody, error) {
	//##########################################
	// Define Container Ports
	//##########################################
	portMaps, portTranslateErr := cg.translatePortsToMaps(templateConfig.Ports)
	if portTranslateErr != nil {
		return nil, fmt.Errorf(
			"container group build create request failed; %s",
			portTranslateErr.Error(),
		)
	}
	//##########################################
	// Define Operating System Type
	//##########################################
	osType, osTypeErr := cg.translateOSType(templateConfig.OperatingSystem)
	if osTypeErr != nil {
		return nil, osTypeErr
	}
	//##########################################
	// Define Container Subnet Id
	//##########################################
	containerGroupSubnet := requests.ContainerGroupSubnetId{}
	containerGroupSubnet.Id = cg.buildSubnetId(
		cg.Subnet.Subscription,
		cg.Subnet.ResourceGroup,
		cg.Subnet.VNetName,
		cg.Subnet.SubnetName,
	)
	containerGroupSubnet.Name = cg.Subnet.SubnetName
	//##########################################
	// Define Container Resources
	//##########################################
	containerResources := requests.Resources{
		ResourceRequest: requests.ResourceRequest{
			CPU:    templateConfig.Resources.CPU,
			Memory: float64(templateConfig.Resources.Memory),
		},
	}
	//##########################################
	// Define Container Instance
	//##########################################
	container := requests.Container{
		Name: templateConfig.ContainerName,
		Properties: requests.ContainerProperties{
			Command:              templateConfig.Commands,
			EnvironmentVariables: cg.translateEnvironmentVars(templateConfig.EnvironmentVariables),
			Image:                templateConfig.Image,
			Ports:                portMaps,
			Resources:            containerResources,
		},
	}
	var containerArr = []requests.Container{}
	containerArr = append(containerArr, container)
	//##########################################
	// Define Create Container Request Body
	//##########################################
	body := requests.CreateContainerGroupBody{
		Location: "canada central",
		Properties: requests.ContainerGroupProperties{
			Containers: containerArr,
			OSType:     osType,
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

func (cg *ContainerGroup) encodeCreateBody(reqBody *requests.CreateContainerGroupBody) (*[]byte, error) {
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
		return fmt.Errorf("failed to create container group")
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")
	_, resErr := cg.azExecuteCreateRequest(request)
	if resErr != nil {
		return resErr
	}
	return nil
}
