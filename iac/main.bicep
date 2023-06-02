@description('Data center location')
param location string

@description('Deployment environment')
param environment string

@description('Virtual network name')
param vnetName string

@description('Name of the container group subnet')
param containerGroupSubnetName string

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-11-01' = {
  name: 'vnet-${vnetName}-${environment}'
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
  }
}

/*
  Creating the virutal network, subnet for required for deploying Azure Container Instances.

*/
resource containerGroupSubnet 'Microsoft.Network/virtualNetworks/subnets@2022-11-01' = {
  name: containerGroupSubnetName
  parent: virtualNetwork
  properties: {
    addressPrefix: '10.0.1.0/24'
    // This delegation is required in order for the Azure Container Instances to deploy
    // inside a virtual network.
    delegations: [
      {
        name: 'Microsoft.ContainerInstance/containerGroups'
        properties: {
          serviceName: 'Microsoft.ContainerInstance/containerGroups'
        }
      }
    ]
  }
  dependsOn: [virtualNetwork]
}
