#!/bin/bash

# Set the options of the getopt command
format=$(getopt -n "$0" -l "subscription:,resource_group:,deployment_location:" -- -- "$@")
if [ $# -lt 3 ]; then
   echo "Wrong number of arguments are passed."
   exit
fi
eval set -- "$format"

# Read the argument values
while [ $# -gt 0 ]
do
     case "$1" in
          --subscription) SUBSCRIPTION="$2"; shift;;
          --resource_group) RESOURCE_GROUP="$2"; shift;;
          --deployment_location) DEPLOYMENT_LOCATION="$2"; shift;;
          --) shift;;
     esac
     shift;
done

echo "Creating resource group [${RESOURCE_GROUP}]"
az group create \
    --name $RESOURCE_GROUP \
    --location $DEPLOYMENT_LOCATION \
    > /dev/null

echo "Creating deployment resources for example architecture"
az deployment group create \
    --resource-group $RESOURCE_GROUP \
    --template-file main.bicep \
    --parameters location=$DEPLOYMENT_LOCATION \
        deploymentEnv='dev' \
        vnetName='data-jobs' \
        containerGroupSubnetName='container-group-subnet' \
        containerGroupSubnetPrefix='10.0.1.0/24' \
        storageAccountName='acitesting' \
        storageAccountSubnetName='storage-account-subnet' \
        storageAccountSubnetPrefix='10.0.2.0/24' \
    > /dev/null

echo "Creating user assigned managed identity."
IDENTITY_JSON=$(az identity create -g $RESOURCE_GROUP -n "example-aci-ua")
IDENTITY_PRINCIPAL_ID=$(jq -r '.principalId' <<< "$IDENTITY_JSON")

echo "Assigning 'Storage Account Contributor' role for managed identity on storage account [blobacitestingdev]"
az role assignment create \
    --assignee-object-id $IDENTITY_PRINCIPAL_ID \
    --assignee-principal-type ServicePrincipal \
    --role 'Storage Account Contributor' \
    --scope /subscriptions/$SUBSCRIPTION/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Storage/storageAccounts/blobacitestingdev \
    > /dev/null