@description('Data center location')
param location string

@description('Deployment environment')
param environment string

@description('Virtual network name')
param vnetName string

@description('Name of the container group subnet')
param containerGroupSubnetName string

resource virutalNetwork 'Microsoft.Network/virtualNetworks@2022-11-01' = {
  name: 'vnet-${vnetName}-${environment}'
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    // subnets: [
    //   {
    //     name: containerGroupSubnetName
    //     properties: {
    //       addressPrefix: '10.0.1.0/24'
    //     }
    //   }
    // ]
  }
}

resource containerGroupSubnet 'Microsoft.Network/virtualNetworks/subnets@2022-11-01' = {
  name: containerGroupSubnetName
  properties: {
    addressPrefix: '10.0.1.0/24'
    delegations: [
      {
        properties: {
          serviceName: 'Microsoft.ContainerInstance/containerGroups'
        }
      }
    ]
  }
}
