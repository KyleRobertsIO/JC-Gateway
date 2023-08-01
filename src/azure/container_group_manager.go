package azure

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	"github.com/kylerobertsio/aci-job-manager/src/api/payloads"
)

type ErrorSource string

const (
	AzureResourceManager ErrorSource = "AZURE_RESOURCE_MANAGER"
	Proxy                ErrorSource = "PROXY"
)

type ContainerGroupManagerError struct {
	Error          string      `json:"error"`
	Source         ErrorSource `json:"source"`
	Code           string      `json:"error_code"`
	HttpStatusCode int         `json:"-"`
}

type ContainerGroupManager struct {
	Credential    *azidentity.DefaultAzureCredential
	Subscription  string
	ResourceGroup string
}

type ContainerGroupState struct {
	ContainerStates []ContainerState `json:"container_states"`
	State           string           `json:"container_group_state"`
}

type ContainerState struct {
	ContainerName string     `json:"container_name"`
	State         string     `json:"state"`
	StartTime     *time.Time `json:"start_time"`
	FinishTime    *time.Time `json:"finish_time"`
	ExitCode      *int32     `json:"exit_code"`
	DetailStatus  string     `json:"detail_status"`
}

func (cgm *ContainerGroupManager) getClient() (*armcontainerinstance.ContainerGroupsClient, *ContainerGroupManagerError) {
	clientFactory, clientFactoryErr := armcontainerinstance.NewClientFactory(
		cgm.Subscription,
		cgm.Credential,
		nil,
	)
	if clientFactoryErr != nil {
		return nil, &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to create client; %s", clientFactoryErr.Error()),
			Source:         Proxy,
			Code:           "ContainerGroupsClientCreateError",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}
	return clientFactory.NewContainerGroupsClient(), nil
}

func (cgm *ContainerGroupManager) CreateOrUpdate(
	payload *payloads.CreateContainerGroup,
) *ContainerGroupManagerError {
	factory := ContainerGroupFactory{Payload: payload}
	containerGroup := factory.Create()
	client, clientErr := cgm.getClient()
	if clientErr != nil {
		return clientErr
	}
	ctx := context.Background()
	_, commandErr := client.BeginCreateOrUpdate(
		ctx,
		cgm.ResourceGroup,
		payload.ContainerGroupName,
		containerGroup,
		nil,
	)
	if commandErr != nil {
		return &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to create/update container group; %s", commandErr.Error()),
			Source:         Proxy,
			Code:           "ContainerGroupCreationError",
			HttpStatusCode: http.StatusBadRequest,
		}
	}
	return nil
}

func (cgm *ContainerGroupManager) determineGroupState(containerStates []ContainerState) string {
	state := "PENDING"
	waiting := 0
	running := 0
	complete := 0
	for _, containerState := range containerStates {
		cState := strings.ToUpper(containerState.State)
		if cState == "TERMINATED" {
			if strings.ToUpper(containerState.DetailStatus) == "ERROR" {
				state = "FAILED"
				break
			} else {
				complete += 1
			}
		}
		if cState == "WAITING" {
			waiting += 1
		}
		if cState == "RUNNING" {
			running += 1
		}
	}
	if waiting > 0 {
		return "WAITING"
	}
	if running > 0 {
		return "RUNNING"
	}
	if complete == len(containerStates) {
		return "COMPLETE"
	}
	return state
}

func (cgm *ContainerGroupManager) Status(
	groupName string,
) (*ContainerGroupState, *ContainerGroupManagerError) {
	client, clientErr := cgm.getClient()
	if clientErr != nil {
		return nil, clientErr
	}
	ctx := context.Background()
	res, commandErr := client.Get(ctx, cgm.ResourceGroup, groupName, nil)
	var resErr *azcore.ResponseError
	if errors.As(commandErr, &resErr) {
		return nil, &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to get container group details; %s", resErr.Error()),
			Source:         AzureResourceManager,
			Code:           resErr.ErrorCode,
			HttpStatusCode: resErr.StatusCode,
		}
	}
	if commandErr != nil {
		return nil, &ContainerGroupManagerError{
			Error:          fmt.Sprintf("failed to get container group details; %s", commandErr.Error()),
			Source:         Proxy,
			Code:           "GetContainerGroupError",
			HttpStatusCode: http.StatusInternalServerError,
		}
	}

	containerStates := make([]ContainerState, 0)
	for _, container := range res.Properties.Containers {
		if container.Properties.InstanceView != nil {
			state := ContainerState{
				ContainerName: *container.Name,
				State:         *container.Properties.InstanceView.CurrentState.State,
				StartTime:     container.Properties.InstanceView.CurrentState.StartTime,
				FinishTime:    container.Properties.InstanceView.CurrentState.FinishTime,
				ExitCode:      container.Properties.InstanceView.CurrentState.ExitCode,
				DetailStatus:  *container.Properties.InstanceView.CurrentState.DetailStatus,
			}
			containerStates = append(containerStates, state)
		}
	}
	groupState := cgm.determineGroupState(containerStates)
	containerGroupState := ContainerGroupState{
		ContainerStates: containerStates,
		State:           groupState,
	}
	return &containerGroupState, nil
}
