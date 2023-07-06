[CmdletBinding()]
Param (
    [Parameter(Mandatory=$true)]
    [string]$Subscription,
    [Parameter(Mandatory=$true)]
    [string]$ResourceGroup,
    [Parameter(Mandatory=$true)]
    [string]$DeploymentLocation
)


Write-Output "Deploying to Azure:"
Write-Output ("Subscription = [{0}] | Resource Group = [{1}] | Location = [{2}]" -f $Subscription, $ResourceGroup, $DeploymentLocation)

Write-Output ('Creating resource group "{0}"' -f $ResourceGroup)
az group create `
    --name $ResourceGroup `
    --location $DeploymentLocation

Write-Output ('Creating example deployment resources from Bicep for "{0}"' -f $ResourceGroup)
az deployment group create `
    --resource-group $ResourceGroup `
    --template-file main.bicep `
    --parameters location=$DeploymentLocation `
        deploymentEnv='dev' `
        vnetName='data-jobs' `
        containerGroupSubnetName='container-group-subnet' `
        containerGroupSubnetPrefix='10.0.1.0/24' `
        storageAccountName='acitesting' `
        storageAccountSubnetName='storage-account-subnet' `
        storageAccountSubnetPrefix='10.0.2.0/24'

Write-Output "Creating user assigned managed identity"
$identity = az identity create `
    -g $ResourceGroup `
    -n "example-aci-ua" | `
    ConvertFrom-Json

Write-Output "Granting IAM [Storage Account Contributor] to example storage account"
az role assignment create `
    --assignee-object-id $identity.principalId `
    --assignee-principal-type ServicePrincipal `
    --role 'Storage Account Contributor' `
    --scope /subscriptions/$Subscription/resourceGroups/$ResourceGroup/providers/Microsoft.Storage/storageAccounts/blobacitestingdev `
    