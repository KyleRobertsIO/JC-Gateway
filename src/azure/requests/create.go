package requests

//##############################################
// General Details
//##############################################

type Port struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

//##############################################
// Container Specific Details
//##############################################

type ResourceRequest struct {
	CPU    int     `json:"cpu"`
	Memory float64 `json:"memoryInGB"`
	// Will need to work in GPU allocation at some point
}

type Resources struct {
	ResourceRequest ResourceRequest `json:"requests"`
}

type ContainerProperties struct {
	Command              []string                 `json:"command"`
	EnvironmentVariables []map[string]interface{} `json:"environmentVariables"`
	Image                string                   `json:"image"`
	Ports                []Port                   `json:"ports"`
	Resources            Resources                `json:"resources"`
}

type Container struct {
	Name       string              `json:"name"`
	Properties ContainerProperties `json:"properties"`
}

//##############################################
// Container Group Specific Details
//##############################################

type ContainerGroupSubnetId struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type IPAddress struct {
	Type  string `json:"type" validate:"omitempty,oneof=Public Private"`
	Ports []Port `json:"ports"`
}

type ContainerGroupProperties struct {
	Containers             []Container              `json:"containers"`
	OSType                 string                   `json:"osType"`
	RestartPolicy          string                   `json:"restartPolicy"`
	ContainerGroupSubnetId []ContainerGroupSubnetId `json:"subnetIds"`
	IPAddress              IPAddress                `json:"ipaddress"`
}

type CreateContainerGroupBody struct {
	Location   string                   `json:"location"`
	Properties ContainerGroupProperties `json:"properties"`
	Name       string                   `json:"name"`
}
