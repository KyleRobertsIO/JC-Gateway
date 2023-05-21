package payloads

import "fmt"

type Subnet struct {
	Subscription       string `json:"subscription"`
	ResourceGroup      string `json:"resource_group"`
	VirtualNetworkName string `json:"virutal_network_name"`
	SubnetName         string `json:"subnet_name"`
}

func (s *Subnet) GetId() string {
	return fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s",
		s.Subscription,
		s.ResourceGroup,
		s.VirtualNetworkName,
		s.SubnetName,
	)
}

type IPAddress struct {
	Type  string `json:"type" validate:"omitempty,oneof=Public Private"`
	Ports []Port `json:"ports"`
}

type EnvironmentVariable struct {
	Secure bool   `json:"secure"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

type Port struct {
	Protocol string `json:"protocol" validate:"omitempty,oneof=tcp udp"`
	Number   int    `json:"number"`
}

type Resources struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
}

type Container struct {
	Name                 string                `json:"name"`
	Ports                []Port                `json:"ports"`
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
	Resources            Resources             `json:"resources"`
	Command              []string              `json:"command"`
	Image                string                `json:"image"`
}

type CreateContainerGroup struct {
	Subscription       string      `json:"subscription"`
	ResourceGroup      string      `json:"resource_group"`
	ContainerGroupName string      `json:"container_group_name"`
	OSType             string      `json:"os_type" default:"Linux" validate:"omitempty,oneof=Linux Windows"`
	Subnet             Subnet      `json:"container_subnet"`
	Containers         []Container `json:"containers"`
	Location           string      `json:"location"`
	IPAddress          IPAddress   `json:"ipaddress"`
}
