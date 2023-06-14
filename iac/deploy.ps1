$Subscription = $args[0]
$ResourceGroupName = "rg-container-groups"
$DeploymentLocation = "canadacentral"

az group create `
    --name $ResourceGroupName `
    --location $DeploymentLocation

az deployment group create `
    --resource-group $ResourceGroupName `
    --template-file main.bicep `
    --parameters location=$DeploymentLocation `
        deploymentEnv='dev' `
        vnetName='data-jobs' `
        containerGroupSubnetName='container-group-subnet' `
        containerGroupSubnetPrefix='10.0.1.0/24' `
        storageAccountName='acitesting' `
        storageAccountSubnetName='storage-account-subnet' `
        storageAccountSubnetPrefix='10.0.2.0/24'

$identity = az identity create `
    -g $ResourceGroupName `
    -n "example-aci-ua" | `
    ConvertFrom-Json

Write-Output $identity

az role assignment create `
    --assignee-object-id $identity.principalId `
    --assignee-principal-type ServicePrincipal `
    --role 'Storage Account Contributor' `
    --scope /subscriptions/$Subscription/resourceGroups/$ResourceGroupName/providers/Microsoft.Storage/storageAccounts/blobacitestingdev
    