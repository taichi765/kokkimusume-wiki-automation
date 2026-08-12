#!/bin/sh
RESOURCE_GROUP_NAME=kokkimusume-discordbot-tfstate-resources
STORAGE_ACCOUNT_NAME=kdbaccount
STORAGE_CONTAINER_NAME=tfstate-storage-container

az group create --name $RESOURCE_GROUP_NAME --location japaneast
az storage account create -g $RESOURCE_GROUP_NAME -n $STORAGE_ACCOUNT_NAME\
    --sku Standard_LRS --encryption-services blob

az storage container create -n $STORAGE_CONTAINER_NAME --account-name $STORAGE_ACCOUNT_NAME