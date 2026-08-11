provider "azurerm" {
  subscription_id = "ad53c5b5-a699-4d4e-8125-0619976eb022"
  features {
    
  }
}

resource "azurerm_resource_group" "kokkimusume-discordbot" {
  name     = "kokkimusume-discordbot-resources"
  location = "Japan East"
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
      image  = "kokkimusumediscordbot-a7g8a0bnajh6fghc.azurecr.io/kokkimusumediscordbot:v0.0.0"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

resource "azurerm_container_registry" "kokkimusume-discordbot" {
  name = "kokkimusumeDiscordbotRegistry"
  resource_group_name = azurerm_resource_group.kokkimusume-discordbot.name
  location = azurerm_resource_group.kokkimusume-discordbot.location
  sku = "Standard"
  admin_enabled = false
}