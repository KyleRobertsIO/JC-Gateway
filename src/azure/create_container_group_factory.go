package azure

import (
	"kyleroberts.io/src/api/payloads"
	"kyleroberts.io/src/azure/requests"
)

type ContainerGroupFactory struct {
	payload payloads.CreateContainerGroup
}

func (cgf *ContainerGroupFactory) Create() requests.CreateContainerGroupBody {
	// Build Container Collection Details
	var containers []requests.Container
	for _, c := range cgf.payload.Containers {
		container := cgf.createContainer(c)
		containers = append(containers, container)
	}
	// Container Group Networking Details
	ipaddress := requests.IPAddress{
		Type:  cgf.payload.IPAddress.Type,
		Ports: cgf.translatePorts(cgf.payload.IPAddress.Ports),
	}
	var subnetCollection []requests.ContainerGroupSubnetId
	if cgf.payload.IPAddress.Type == "Private" {
		subnetId := requests.ContainerGroupSubnetId{
			Id:   cgf.payload.Subnet.GetId(),
			Name: cgf.payload.Subnet.SubnetName,
		}
		subnetCollection = append(subnetCollection, subnetId)
	}
	// Build Container Group Details
	groupProps := requests.ContainerGroupProperties{
		Containers:             containers,
		OSType:                 cgf.payload.OSType,
		ContainerGroupSubnetId: subnetCollection,
		IPAddress:              ipaddress,
	}
	return requests.CreateContainerGroupBody{
		Location:   cgf.payload.Location,
		Properties: groupProps,
	}
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

func (cgf *ContainerGroupFactory) createContainer(
	container payloads.Container,
) requests.Container {
	containerResources := requests.Resources{
		ResourceRequest: requests.ResourceRequest{
			CPU:    container.Resources.CPU,
			Memory: float64(container.Resources.Memory),
		},
	}
	containerProps := requests.ContainerProperties{
		Command:              container.Command,
		EnvironmentVariables: cgf.translateEnvVars(container.EnvironmentVariables),
		Image:                container.Image,
		Ports:                cgf.translatePorts(container.Ports),
		Resources:            containerResources,
	}
	return requests.Container{
		Name:       container.Name,
		Properties: containerProps,
	}
}
