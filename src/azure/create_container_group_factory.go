package azure

import (
	"kyleroberts.io/src/api/payloads"
	"kyleroberts.io/src/azure/requests"
)

type ContainerGroupFactory struct {
	payloads.CreateContainerGroup
}

func (cgf *ContainerGroupFactory) Create() requests.CreateContainerGroupBody {

}

func (cgf *ContainerGroupFactory) translateEnvVars(
	envVars []payloads.EnvironmentVariable,
) []map[string]interface{} {
	var envVarCollection []map[string]interface{}
	for _, port := range envVars {
		translatedEnvVar := cgf.translateEnvVar(port)
		envVarCollection = append(envVarCollection, translatedEnvVar)
	}
	return envVarCollection
}

func (cgf *ContainerGroupFactory) translateEnvVar(
	envVar payloads.EnvironmentVariable,
) map[string]interface{} {
	translatedEnvVar := make(map[string]interface{})
	translatedEnvVar["name"] = envVar.Name
	if envVar.Secure {
		translatedEnvVar["secureValue"] = envVar.Value
	} else {
		translatedEnvVar["value"] = envVar.Value
	}
	return translatedEnvVar
}

func (cgf *ContainerGroupFactory) translatePorts(ports []payloads.Port) []requests.Port {
	var translatedPortCollection []requests.Port
	for _, port := range ports {
		translatedPort := cgf.translatePort(port)
		translatedPortCollection = append(translatedPortCollection, translatedPort)
	}
	return translatedPortCollection
}

func (cgf *ContainerGroupFactory) translatePort(port payloads.Port) requests.Port {
	translatedPort := requests.Port{
		Port:     port.Number,
		Protocol: port.Protocol,
	}
	return translatedPort
}

func (cgf *ContainerGroupFactory) createContainerGroup(
	payload payloads.Container,
	osType string,
) requests.Container {
	containerResources := requests.Resources{
		ResourceRequest: requests.ResourceRequest{
			CPU:    payload.Resources.CPU,
			Memory: float64(payload.Resources.Memory),
		},
	}
	containerProps := requests.ContainerProperties{
		Command:              payload.Command,
		EnvironmentVariables: cgf.translateEnvVars(payload.EnvironmentVariables),
		Image:                payload.Image,
		Ports:                cgf.translatePorts(payload.Ports),
		Resources:            containerResources,
	}
	container := requests.Container{
		Name:       payload.Name,
		Properties: containerProps,
	}
	return container
}
