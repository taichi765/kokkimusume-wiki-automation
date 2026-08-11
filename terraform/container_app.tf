

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