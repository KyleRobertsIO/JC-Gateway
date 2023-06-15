@description('Data center location')
param location string

@description('Deployment environment')
param deploymentEnv string

@description('Virtual network name')
param vnetName string

@description('Name of the container group subnet')
param containerGroupSubnetName string
@description('Azure container instance subnet address prefix')
param containerGroupSubnetPrefix string

@description('Azure storage account name')
param storageAccountName string
@description('Name of storage account subnet')
param storageAccountSubnetName string
@description('Azure storage account subnet address prefix')
param storageAccountSubnetPrefix string

var blobStoragePrivateDnsZoneName = 'privatelink.blob.${environment().suffixes.storage}'

/*
  Define sample virtual network for the deployment of the architecture.
*/
resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-11-01' = {
  name: 'vnet-${vnetName}-${deploymentEnv}'
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: ['10.0.0.0/16']
    }
  }
}

/*
  Creating the subnet required for deploying Azure Container Instances.
*/
resource containerGroupSubnet 'Microsoft.Network/virtualNetworks/subnets@2022-11-01' = {
  name: containerGroupSubnetName
  parent: virtualNetwork
  properties: {
    addressPrefix: containerGroupSubnetPrefix
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
}

/*
  Creating the subnet for general resources like Azure Storage Accounts.
*/
resource storageAccountSubnet 'Microsoft.Network/virtualNetworks/subnets@2022-11-01' = {
  name: storageAccountSubnetName
  parent: virtualNetwork
  properties: {
    addressPrefix: storageAccountSubnetPrefix
  }
}

/*
  Create a Azure Storage Account to demo connecting to resources inside
  of a virtual network with Azure Container Instances.
*/
resource storageAccount 'Microsoft.Storage/storageAccounts@2022-09-01' = {
  name: 'blob${storageAccountName}${deploymentEnv}'
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'BlobStorage'
  properties: {
    accessTier: 'Hot'
    allowBlobPublicAccess: false
    allowCrossTenantReplication: false
    allowedCopyScope: 'PrivateLink'
    allowSharedKeyAccess: true
    defaultToOAuthAuthentication: false
    dnsEndpointType: 'Standard'
    immutableStorageWithVersioning: {
      enabled: false
    }
    isHnsEnabled: false
    isLocalUserEnabled: false
    isNfsV3Enabled: false
    isSftpEnabled: false
    keyPolicy: {
      keyExpirationPeriodInDays: 90
    }
    minimumTlsVersion: 'TLS1_2'
    publicNetworkAccess: 'Disabled'
    routingPreference: {
      publishInternetEndpoints: false
      publishMicrosoftEndpoints: true
      routingChoice: 'MicrosoftRouting'
    }
    supportsHttpsTrafficOnly: true
  }
}

/*
  Create a Private Endpoint attached to the create Storage Account for private
  virtual network traffic through the Azure data center.
*/
resource blobStoragePrivateEndpoint 'Microsoft.Network/privateEndpoints@2022-11-01' = {
  name: 'nic-blob-${storageAccountName}'
  location: location
  properties: {
    subnet: {
      id: storageAccountSubnet.id
    }
    privateLinkServiceConnections: [
      {
        name: 'nic-blob-${storageAccountName}'
        properties: {
          privateLinkServiceId: storageAccount.id
          groupIds: [
            'blob'
          ]
          privateLinkServiceConnectionState: {
            status: 'Approved'
            description: 'Auto-Approved'
            actionsRequired: 'None'
          }
        }
      }
    ]
  }
}

/*
  Create a Private DNS Zone for blob storage.
*/
resource blobPrivateDnsZone 'Microsoft.Network/privateDnsZones@2020-06-01' = {
  name: blobStoragePrivateDnsZoneName
  location: 'global'
  properties: {}
  dependsOn: [blobStoragePrivateEndpoint]
}

resource privateEndpointDns 'Microsoft.Network/privateEndpoints/privateDnsZoneGroups@2022-01-01' = {
  parent: blobStoragePrivateEndpoint
  name: 'blob-PrivateDnsZoneGroup'
  properties:{
    privateDnsZoneConfigs: [
      {
        name: blobStoragePrivateDnsZoneName
        properties:{
          privateDnsZoneId: blobPrivateDnsZone.id
        }
      }
    ]
  }
  dependsOn: [storageAccount]
}

resource privateDnsZoneLink 'Microsoft.Network/privateDnsZones/virtualNetworkLinks@2020-06-01' = {
  parent: blobPrivateDnsZone
  name: '${blobStoragePrivateDnsZoneName}-link'
  location: 'global'
  properties: {
    registrationEnabled: true
    virtualNetwork: {
      id: virtualNetwork.id
    }
  }
}
