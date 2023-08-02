package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	"github.com/kylerobertsio/aci-job-manager/src/api/payloads"
	"github.com/kylerobertsio/aci-job-manager/src/azure"
)

func CreateContainerGroup() armcontainerinstance.ContainerGroup {
	containerPorts := make([]payloads.Port, 0)
	port := payloads.Port{Protocol: "tcp", Number: 443}
	containerPorts = append(containerPorts, port)

	containerCollection := make([]payloads.Container, 0)
	container := payloads.Container{
		Name:                 "example-container-name",
		Ports:                containerPorts,
		EnvironmentVariables: make([]payloads.EnvironmentVariable, 0),
		Command:              nil,
		Image:                "nginx",
		Resources: payloads.Resources{
			CPU:    1.0,
			Memory: 1.0,
		},
	}
	containerCollection = append(containerCollection, container)

	payload := payloads.CreateContainerGroup{
		Subscription:       "234234214-2342432dsf2-sdf21331",
		ResourceGroup:      "rg-autotest",
		ContainerGroupName: "example-group-name",
		OSType:             "Linux",
		Location:           "canadacentral",
		IPAddress: payloads.IPAddress{
			Type:  "Private",
			Ports: make([]payloads.Port, 0),
		},
		Subnet: payloads.Subnet{
			Subscription:       "234234214-2342432dsf2-sdf21331",
			ResourceGroup:      "rg-autotest",
			VirtualNetworkName: "vnet-example-name",
			SubnetName:         "subnet-example-name",
		},
		Containers: containerCollection,
	}

	factory := azure.ContainerGroupFactory{Payload: &payload}
	cg := factory.Create()
	return cg
}

func TestTranslatedContainerPorts(t *testing.T) {
	cg := CreateContainerGroup()
	translatedPortProtocol := *cg.Properties.Containers[0].Properties.Ports[0].Protocol
	if translatedPortProtocol != "tcp" {
		t.Error("expected container port to be of Protocol 'tcp'")
	}

	translatedPortNumber := *cg.Properties.Containers[0].Properties.Ports[0].Port
	if translatedPortNumber != 443 {
		t.Error("expected container port to be of Number 443")
	}
}

func TestTranslatedLocation(t *testing.T) {
	cg := CreateContainerGroup()
	translatedLocation := *cg.Location
	expectedLocation := "canadacentral"
	if translatedLocation != expectedLocation {
		t.Error("container group doesn't container Location = 'canadacentral'")
	}
}
