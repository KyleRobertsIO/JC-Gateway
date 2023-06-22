using 'main.bicep'

param location = 'canada central'
param deploymentEnv = 'dev'

// Virtual Network
param vnetName = 'data-jobs'

// Container Instance Network
param containerGroupSubnetName = 'container-group-subnet'
param containerGroupSubnetPrefix = '10.0.1.0/24'

// Storage Account
param storageAccountName = 'acitesting'

// Storage Account Network
param storageAccountSubnetName = 'storage-account-subnet'
param storageAccountSubnetPrefix = '10.0.2.0/24'
