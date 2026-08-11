

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

  identity {
    type = "SystemAssigned"
  }

  registry {
    server = azurerm_container_registry.acr.login_server
    identity = "system"
  }

  ingress {
    target_port = 8080
    external_enabled = true
    allow_insecure_connections = false
   
    traffic_weight {
      percentage = 100
      latest_revision = true
    }
  }

  template {
    min_replicas = 0
    max_replicas = 1

    container {
      name   = "kokkimusume-discordbot-container"
      image  =  "${azurerm_container_registry.acr.login_server}/kokkimusumediscordbot:v0.0.0"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

resource "azurerm_role_assignment" "acr_pull" {
  scope = azurerm_container_registry.acr.id
  role_definition_name = "AcrPull"
  principal_id = azurerm_container_app.kokkimusume-discordbot.identity[0].principal_id
}