#!/bin/bash

# Log into Azure
az login

# Show current subscription (use 'az account set' to change subscription)
az account show

# variables
subscriptionId=$(az account show --query id -o tsv)
appName="gogoto-app_sp"
RBACRole="Contributor"

githubOrgName="MoreiraVinicius"
githubRepoName="gogoto"
githubBranch="feature/github-actions"

# Create AAD App and Principal
appId=$(az ad app create --display-name $appName --query appId -o tsv)
az ad sp create --id $appId

# Create federated GitHub credentials (Entity type 'Branch')
githubBranchConfig='{
    "name": "GH-['$githubOrgName'-'$githubRepoName']-Branch-['$githubBranch']",
    "issuer": "https://token.actions.githubusercontent.com",
    "subject": "repo:'$githubOrgName'/'$githubRepoName':ref:refs/heads/'$githubBranch'",
    "description": "Federated credential linked to GitHub ['$githubBranch'] branch @: ['$githubOrgName'/'$githubRepoName']",
    "audiences": ["api://AzureADTokenExchange"]
}'
echo $githubBranchConfig | az ad app federated-credential create --id $appId --parameters @-

# Create federated GitHub credentials (Entity type 'Pull Request')
githubPRConfig='{
    "name": "GH-['$githubOrgName'-'$githubRepoName']-PR",
    "issuer": "https://token.actions.githubusercontent.com",
    "subject": "repo:'$githubOrgName'/'$githubRepoName':pull_request",
    "description": "Federated credential linked to GitHub Pull Requests @: ['$githubOrgName'/'$githubRepoName']",
    "audiences": ["api://AzureADTokenExchange"]
}'
echo $githubPRConfig | az ad app federated-credential create --id $appId --parameters @-

# Assign RBAC permissions to Service Principal (Change as necessary)
az role assignment create --role $RBACRole --assignee $appId --subscription $subscriptionId

# Permission 2 (Example)
# az role assignment create --role "Reader and Data Access" --assignee $appId --scope "/subscriptions/$subscriptionId/resourceGroups/$resourceGroupName/providers/Microsoft.Storage/storageAccounts/$storageName"