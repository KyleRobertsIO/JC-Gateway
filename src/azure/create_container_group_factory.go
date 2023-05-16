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

func (cgf *ContainerGroupFactory) createContainerGroup() requests.Container {

}
