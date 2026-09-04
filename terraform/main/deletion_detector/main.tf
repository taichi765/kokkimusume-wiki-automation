
data "azurerm_resource_group" "detector" {
  name = "deletion-detector"
}

data "azurerm_client_config" "current" {

}

resource "azurerm_key_vault" "detector" {
  name                = "kokkimusume-detector"
  location            = data.azurerm_resource_group.detector.location
  resource_group_name = data.azurerm_resource_group.detector.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                   = "standard"
  rbac_authorization_enabled = true
}


resource "azurerm_log_analytics_workspace" "detector" {
  name                = "deletion-detector-logs"
  location            = data.azurerm_resource_group.detector.location
  resource_group_name = data.azurerm_resource_group.detector.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_job" "detector" {
  name                         = "deletion-detector"
  location                     = data.azurerm_resource_group.detector.location
  resource_group_name          = data.azurerm_resource_group.detector.name
  container_app_environment_id = var.container_app_environment_id

  replica_timeout_in_seconds = 60
  replica_retry_limit        = 5

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.detector.id]
  }

  registry {
    server   = var.acr_login_server
    identity = azurerm_user_assigned_identity.detector.id
  }

  secret {
    name                = "discord-token"
    key_vault_secret_id = "${azurerm_key_vault.detector.vault_uri}secrets/discord-token"
    identity            = azurerm_user_assigned_identity.detector.id
  }

  secret {
    name                = "wikiwiki-password"
    key_vault_secret_id = "${azurerm_key_vault.detector.vault_uri}secrets/wikiwiki-password"
    identity            = azurerm_user_assigned_identity.detector.id
  }

  #schedule_trigger_config {
    # 一時間に一回
  #  cron_expression = "0 */1 * * *"
  #}

  manual_trigger_config {
    
  }

  template {
    container {
      image = "${var.acr_login_server}/deletion-detector:${var.commit_hash}"
      name  = "deletion-detector"
      readiness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/readiness"
      }

      startup_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/startup"
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/liveness"
      }

      env {
        name        = "DISCORD_TOKEN"
        secret_name = "discord-token"
      }

      env {
        name        = "WIKIWIKI_PASSWORD"
        secret_name = "wikiwiki-password"
      }

      env {
        name = "AZURE_CLIENT_ID"
        value = azurerm_user_assigned_identity.detector.client_id
      }

      cpu    = 0.5
      memory = "1Gi"
    }
  }
}

