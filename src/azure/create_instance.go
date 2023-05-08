package azure

import (
	"fmt"
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
	Name       string
	Properties ContainerProperties `json:"properties"`
}

type RequestProperties struct {
	Containers []Container
}

type CreateRequestBody struct {
	Location   string
	Properties RequestProperties
}

//##############################
// Inbound Client Payload
//##############################

type ContainerGroup struct {
	Subscription  string
	ResourceGroup string
	Name          string
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

func (cg *ContainerGroup) Create(
	apiVersion string,
	accessToken string,
	templateConfig templates.YamlConfig,
) error {
	client := &http.Client{}
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
	request, err := http.NewRequest(
		http.MethodPut,
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to create container group.")
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, reqErr := client.Do(request)
	if reqErr != nil {
		return fmt.Errorf("HTTP request to %s failed; %s", url, reqErr.Error())
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(bodyBytes))
	return nil
}
