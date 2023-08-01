package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	"github.com/kylerobertsio/aci-job-manager/src/api/payloads"
)

type ContainerGroupFactory struct {
	Payload *payloads.CreateContainerGroup
}

func (cgf *ContainerGroupFactory) translateContainerPort(
	port payloads.Port,
) armcontainerinstance.ContainerPort {
	return armcontainerinstance.ContainerPort{
		Port:     (*int32)(&port.Number),
		Protocol: (*armcontainerinstance.ContainerNetworkProtocol)(&port.Protocol),
	}
}

func (cgf *ContainerGroupFactory) translateContainerPorts(
	ports []payloads.Port,
) []*armcontainerinstance.ContainerPort {
	translatedPortCollection := make([]*armcontainerinstance.ContainerPort, 0)
	for _, port := range ports {
		translatedPort := cgf.translateContainerPort(port)
		translatedPortCollection = append(translatedPortCollection, &translatedPort)
	}
	return translatedPortCollection
}

func (cgf *ContainerGroupFactory) translateEnvVar(
	envVar payloads.EnvironmentVariable,
) armcontainerinstance.EnvironmentVariable {
	if envVar.Secure {
		return armcontainerinstance.EnvironmentVariable{
			Name:        &envVar.Name,
			SecureValue: &envVar.Value,
		}
	} else {
		return armcontainerinstance.EnvironmentVariable{
			Name:  &envVar.Name,
			Value: &envVar.Value,
		}
	}
}

func (cgf *ContainerGroupFactory) translateEnvVars(
	envVars []payloads.EnvironmentVariable,
) []*armcontainerinstance.EnvironmentVariable {
	envVarCollection := make([]*armcontainerinstance.EnvironmentVariable, 0)
	for _, port := range envVars {
		translatedEnvVar := cgf.translateEnvVar(port)
		envVarCollection = append(envVarCollection, &translatedEnvVar)
	}
	return envVarCollection
}

func (cgf *ContainerGroupFactory) createContainer(
	container payloads.Container,
) armcontainerinstance.Container {
	resourceRequirements := armcontainerinstance.ResourceRequirements{
		Requests: &armcontainerinstance.ResourceRequests{
			CPU:        &container.Resources.CPU,
			MemoryInGB: &container.Resources.Memory,
		},
	}
	envVars := cgf.translateEnvVars(container.EnvironmentVariables)
	ports := cgf.translateContainerPorts(container.Ports)
	properties := armcontainerinstance.ContainerProperties{
		Image:                &container.Image,
		Resources:            &resourceRequirements,
		Command:              container.Command,
		EnvironmentVariables: envVars,
		Ports:                ports,
	}
	return armcontainerinstance.Container{
		Name:       &container.Name,
		Properties: &properties,
	}
}

func (cgf *ContainerGroupFactory) translatePort(
	port payloads.Port,
) armcontainerinstance.Port {
	return armcontainerinstance.Port{
		Port:     &port.Number,
		Protocol: (*armcontainerinstance.ContainerGroupNetworkProtocol)(&port.Protocol),
	}
}

func (cgf *ContainerGroupFactory) translatePorts(
	ports []payloads.Port,
) []*armcontainerinstance.Port {
	translatedPortCollection := make([]*armcontainerinstance.Port, 0)
	for _, port := range ports {
		translatedPort := cgf.translatePort(port)
		translatedPortCollection = append(translatedPortCollection, &translatedPort)
	}
	return translatedPortCollection
}

func (cgf *ContainerGroupFactory) translateIPAddress(
	ipAddress *payloads.IPAddress,
) *armcontainerinstance.IPAddress {
	ports := cgf.translatePorts(ipAddress.Ports)
	return &armcontainerinstance.IPAddress{
		Type:  (*armcontainerinstance.ContainerGroupIPAddressType)(&ipAddress.Type),
		Ports: ports,
	}
}

func (cgf *ContainerGroupFactory) Create() armcontainerinstance.ContainerGroup {
	// Create the groups containers
	containerCollection := make([]*armcontainerinstance.Container, 0)
	for _, container := range cgf.Payload.Containers {
		container := cgf.createContainer(container)
		containerCollection = append(containerCollection, &container)
	}
	// IP Address
	ipAddress := cgf.translateIPAddress(&cgf.Payload.IPAddress)
	subnetIDs := make([]*armcontainerinstance.ContainerGroupSubnetID, 0)
	if armcontainerinstance.ContainerGroupIPAddressType(cgf.Payload.IPAddress.Type) == armcontainerinstance.ContainerGroupIPAddressTypePrivate {
		subnetIDStr := cgf.Payload.Subnet.GetId()
		subnetID := armcontainerinstance.ContainerGroupSubnetID{
			ID:   &subnetIDStr,
			Name: &cgf.Payload.Subnet.SubnetName,
		}
		subnetIDs = append(subnetIDs, &subnetID)
	}
	// Restart Policy
	restartPolicy := armcontainerinstance.ContainerGroupRestartPolicyNever
	// Create Container Group Properties
	containerGroupProps := armcontainerinstance.ContainerGroupPropertiesProperties{
		Containers:    containerCollection,
		OSType:        (*armcontainerinstance.OperatingSystemTypes)(&cgf.Payload.OSType),
		IPAddress:     ipAddress,
		RestartPolicy: &restartPolicy,
		SubnetIDs:     subnetIDs,
	}
	// Create the container group
	return armcontainerinstance.ContainerGroup{
		Location:   &cgf.Payload.Location,
		Name:       &cgf.Payload.ContainerGroupName,
		Properties: &containerGroupProps,
	}
}
