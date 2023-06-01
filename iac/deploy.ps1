$ResourceGroupName = "rg-container-groups"
$DeploymentLocation = "canadacentral"

az group create `
    --name $ResourceGroupName `
    --location $DeploymentLocation

az deployment group create `
    --resource-group $ResourceGroupName `
    --template-file main.bicep `
    --parameters location=$DeploymentLocation `
        environment='dev' `
        vnetName='data-jobs' `
        containerGroupSubnetName='container-group-subnet'