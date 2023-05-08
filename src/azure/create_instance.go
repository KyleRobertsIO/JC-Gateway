package azure

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
	Command              []string           `json:"command"`
	EnvironmentVariables []string           `json:"environmentVariables"`
	Image                string             `json:"image"`
	Ports                []map[string]int   `json:"ports"`
	Resources            ContainerResources `json:"resources"`
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

func (cg *ContainerGroup) translatePorts(ports []int) []map[string]int {
	var translatedPorts []map[string]int
	for index, port := range ports {
		translatedPorts = append(translatedPorts, port)
	}
}

func (cg *ContainerGroup) buildCreateRequestBody(
	templateConfig templates.YamlConfig,
) {
	containerResouces := ContainerResources{
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
			Ports:                templateConfig.Ports,
			Resources:            containerResouces,
		},
	}
	body := CreateRequestBody{
		Location: "canada-central",
		Properties: RequestProperties{
			Containers: [1]Container{container},
		},
	}
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
