provider "azurerm" {
  subscription_id = "ad53c5b5-a699-4d4e-8125-0619976eb022"
  features {
    
  }
}

provider "azuread" {
  
}

resource "azurerm_resource_group" "kokkimusume-discordbot" {
  name     = "kokkimusume-discordbot-resources"
  location = "Japan East"
}

resource "azuread_application" "github-actions" {
  display_name = "kokkimusume-discordbot-github-actions"
}

resource "azuread_service_principal" "github-actions" {
  client_id = azuread_application.github-actions.client_id
}

resource "azuread_application_federated_identity_credential" "github-actions" {
  application_id = azuread_application.github-actions.id
  display_name = "kokkimusume-discordbot-github-actions-cred"
  description = "Github Actions OIDC federation"
  audiences = ["api://AzureADTokenExchange"]
  issuer = "https://tokens.actions.githubusercontent.com"
  subject = "repo:taichi765/kokkimusume-wiki-automation:ref:refs/heads/main"
}

resource "azurerm_role_assignment" "acr_push" {
  scope = azurerm_container_registry.acr.id
  role_definition_name = "AcrPush"
  principal_id = azuread_service_principal.github-actions.object_id
}

resource "azurerm_container_registry" "acr" {
  name = "kokkimusumeDiscordbotRegistry"
  resource_group_name = azurerm_resource_group.kokkimusume-discordbot.name
  location = azurerm_resource_group.kokkimusume-discordbot.location
  sku = "Standard"
  admin_enabled = false
}

resource "azurerm_log_analytics_workspace" "kokkimusume-discordbot" {
  name                = "kokkimusume-discordbot-logs"
  location            = azurerm_resource_group.kokkimusume-discordbot.location
  resource_group_name = azurerm_resource_group.kokkimusume-discordbot.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "kokkimusume-discordbot" {
  name                       = "kokkimusume-discordbot-env"
  location                   = azurerm_resource_group.kokkimusume-discordbot.location
  resource_group_name        = azurerm_resource_group.kokkimusume-discordbot.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.kokkimusume-discordbot.id
  logs_destination = "log-analytics"
}

resource "azurerm_container_app" "kokkimusume-discordbot" {
  name                         = "kokkimusume-discordbot-app"
  container_app_environment_id = azurerm_container_app_environment.kokkimusume-discordbot.id
  resource_group_name          = azurerm_resource_group.kokkimusume-discordbot.name
  revision_mode                = "Single"

  template {
    container {
      name   = "kokkimusume-discordbot-container"
      image  =  "mcr.microsoft.com/k8se/quickstart:latest"#"kokkimusumediscordbot-a7g8a0bnajh6fghc.azurecr.io/kokkimusumediscordbot:v0.0.0"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

data "azurerm_client_config" "current" {
  
}

output "azure_client_id" {
  value = azuread_application.github-actions.client_id
}

output "azure_tenant_id"{
  value = data.azurerm_client_config.current.tenant_id
}

output "azure_subscription_id" {
  value = data.azurerm_client_config.current.subscription_id
}