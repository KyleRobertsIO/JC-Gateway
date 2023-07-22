@description('The SKU of App Service Plan ')
param sku string = 'B1'

@description('Location for all resources.')
param location string = resourceGroup().location

resource appServicePlan 'Microsoft.Web/serverfarms@2022-09-01' = {
  name: 'as-cacn-job-manager-dev'
  location: location
  properties: {
    reserved: true
  }
  sku: {
    name: sku
  }
  kind: 'linux'
}

resource webApp 'Microsoft.Web/sites@2021-01-01' = {
  name: 'asp-cacn-job-manager-dev'
  location: location
  tags: {}
  kind: 'app,linux,container'
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    httpsOnly: true
    siteConfig: {
      alwaysOn: true
      minTlsVersion: '1.2'
      linuxFxVersion: 'DOCKER|scuffedfox/az_job_container_manager:alpha'
      appSettings: [
        {
          name: 'DOCKER_REGISTRY_SERVER_URL'
          value: 'registry.hub.docker.com'
        }
        {
          name: 'WEBSITES_ENABLE_APP_SERVICE_STORAGE'
          value: 'false'
        }
        {
          name: 'WEBSITES_CONTAINER_START_TIME_LIMIT'
          value: '90'
        }
        {
          name: 'WEBSITES_PORT'
          value: '8080'
        }
        {
          name: 'WEBSITE_WARMUP_PATH'
          value: '/api/ping'
        }
        {
          name: 'WEBSITE_SWAP_WARMUP_PING_PATH'
          value: '/api/ping'
        }
        {
          name: 'PORT'
          value: '8080'
        }
        {
          name: 'LOGGER_LOG_LEVEL'
          value: 'INFO'
        }
        {
          name: 'LOGGER_FILE_PATH'
          value: './default.json'
        }
        {
          name: 'GIN_PORT'
          value: '8080'
        }
        {
          name: 'GIN_MODE'
          value: 'release'
        }
        {
          name: 'AZURE_AUTH_TYPE'
          value: 'SERVICE_PRINCIPAL'
        }
        {
          name: 'AZURE_AUTH_CLIENT_ID'
          value: '9e7c69a7-c510-403d-b4a1-59db42b65834'
        }
        {
          name: 'AZURE_AUTH_CLIENT_SECRET'
          value: 'Ho08Q~ikOZk61_WJzu5-qDAI7JwWNjQO-yG-.b56'
        }
        {
          name: 'AZURE_AUTH_TENANT_ID'
          value: '1a68d4be-bfde-4a10-847b-ca21a96d6ef4'
        }
      ]
    }
    serverFarmId: appServicePlan.id 
  }
}
