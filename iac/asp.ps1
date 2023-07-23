[CmdletBinding()]
Param (
    [Parameter(Mandatory=$true)]
    [string]$ResourceGroup,
    [Parameter(Mandatory=$true)]
    [string]$DeploymentLocation
)


az group create `
    --name $ResourceGroup `
    --location $DeploymentLocation

az deployment group create `
    --resource-group $ResourceGroup `
    --template-file app_service.bicep `
    --parameters location=$DeploymentLocation `
        sku='B1'