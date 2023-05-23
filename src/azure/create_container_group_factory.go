package azure

import (
	"kyleroberts.io/src/api/payloads"
	"kyleroberts.io/src/azure/requests"
)

type ContainerGroupFactory struct {
	Payload *payloads.CreateContainerGroup
}

func (cgf *ContainerGroupFactory) Create() requests.CreateContainerGroupBody {
	// Build Container Collection Details
	containers := make([]requests.Container, 0)
	for _, c := range cgf.Payload.Containers {
		container := cgf.createContainer(c)
		containers = append(containers, container)
	}
	// Container Group Networking Details
	ipaddress := requests.IPAddress{
		Type:  cgf.Payload.IPAddress.Type,
		Ports: cgf.translatePorts(cgf.Payload.IPAddress.Ports),
	}
	subnetCollection := make([]requests.ContainerGroupSubnetId, 0)
	if cgf.Payload.IPAddress.Type == "Private" {
		subnetId := requests.ContainerGroupSubnetId{
			Id:   cgf.Payload.Subnet.GetId(),
			Name: cgf.Payload.Subnet.SubnetName,
		}
		subnetCollection = append(subnetCollection, subnetId)
	}
	// Build Container Group Details
	groupProps := requests.ContainerGroupProperties{
		Containers:             containers,
		OSType:                 cgf.Payload.OSType,
		ContainerGroupSubnetId: subnetCollection,
		IPAddress:              ipaddress,
		RestartPolicy:          "Never",
	}
	return requests.CreateContainerGroupBody{
		Location:   cgf.Payload.Location,
		Properties: groupProps,
		Name:       cgf.Payload.ContainerGroupName,
	}
}

func (cgf *ContainerGroupFactory) translateEnvVars(
	envVars []payloads.EnvironmentVariable,
) []map[string]interface{} {
	envVarCollection := make([]map[string]interface{}, 0)
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
	translatedPortCollection := make([]requests.Port, 0)
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
