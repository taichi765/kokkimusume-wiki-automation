# discord_botとdeletion_detectorの間で共有されるリソース
data "azurerm_resource_group" "shared" {
  name = "kokkimusume-shared"
}

resource "azurerm_container_registry" "shared" {
  name                = "kokkimusumeDiscordbotRegistry"
  resource_group_name = data.azurerm_resource_group.shared.name
  location            = data.azurerm_resource_group.shared.location
  sku                 = "Standard"
  admin_enabled       = false
}

resource "azurerm_log_analytics_workspace" "shared" {
  name                = "log-analytics-workspace"
  location            = data.azurerm_resource_group.shared.location
  resource_group_name = data.azurerm_resource_group.shared.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "shared" {
  name                       = "app-env"
  location                   = data.azurerm_resource_group.shared.location
  resource_group_name        = data.azurerm_resource_group.shared.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.shared.id
  logs_destination           = "log-analytics"
  #logs_destination           = "azure-monitor"
}


resource "azurerm_monitor_diagnostic_setting" "containerapp" {
  name                       = "containerapp"
  target_resource_id         = azurerm_container_app_environment.shared.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.shared.id

  enabled_log {
    category = "ContainerAppHTTPLogs"
  }
}


module "discord_bot" {
  source     = "./discord_bot"
  depends_on = [azurerm_container_registry.shared]

  github-installation-id = var.github-installation-id
  github-app-id          = var.github-app-id
  discord-app-id         = var.discord-app-id
  discord-public-key     = var.discord-public-key
  commit_sha256          = var.discord_bot_commit_hash

  acr_group_name               = data.azurerm_resource_group.shared.name
  acr_name                     = azurerm_container_registry.shared.name
  container_app_environment_id = azurerm_container_app_environment.shared.id
}

module "deletion_detector" {
  source = "./deletion_detector"

  acr_id                       = azurerm_container_registry.shared.id
  acr_login_server             = azurerm_container_registry.shared.login_server
  commit_hash                  = var.deletion_detector_commit_hash
  container_app_environment_id = azurerm_container_app_environment.shared.id
}
