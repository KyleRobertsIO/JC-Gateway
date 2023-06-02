#!/bin/bash
RESOURCE_GROUP_NAME="rg-container-groups"       
DEPLOYMENT_LOCATION="canadacentral"

az group create \
    --name $RESOURCE_GROUP_NAME \
    --location $DEPLOYMENT_LOCATION

az deployment group create \
    --resource-group $RESOURCE_GROUP_NAME \
    --template-file main.bicep \
    --parameters location=$DEPLOYMENT_LOCATION \
        environment="dev" \
        vnetName="data-jobs" \
        containerGroupSubnetPrefix="10.0.1.0/24" \
        containerGroupSubnetName="container-group-subnet"