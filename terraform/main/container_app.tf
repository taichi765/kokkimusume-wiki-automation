resource "azurerm_log_analytics_workspace" "app" {
  name                = "log-analytics-workspace"
  location            = data.azurerm_resource_group.app.location
  resource_group_name = data.azurerm_resource_group.app.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_key_vault" "app" {
  name                = "kdb-key-vault"
  location            = data.azurerm_resource_group.app.location
  resource_group_name = data.azurerm_resource_group.app.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                   = "standard"
  rbac_authorization_enabled = true
}

resource "azurerm_container_app_environment" "app" {
  name                       = "app-env"
  location                   = data.azurerm_resource_group.app.location
  resource_group_name        = data.azurerm_resource_group.app.name
  logs_destination           = "azure-monitor"
}

resource "azurerm_container_app" "app" {
  name                         = "app"
  container_app_environment_id = azurerm_container_app_environment.app.id
  resource_group_name          = data.azurerm_resource_group.app.name
  revision_mode                = "Single"

  depends_on = [azurerm_role_assignment.acr_pull, azurerm_role_assignment.key_vault]

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.app.id]
  }

  registry {
    server   = azurerm_container_registry.acr.login_server
    identity = azurerm_user_assigned_identity.app.id
  }

  ingress {
    target_port                = 8080
    external_enabled           = true
    allow_insecure_connections = false

    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }

  secret {
    name                = "discord-token"
    key_vault_secret_id = "${azurerm_key_vault.app.vault_uri}secrets/discord-token"
    identity            = azurerm_user_assigned_identity.app.id
  }

  secret {
    name                = "github-private-key"
    key_vault_secret_id = "${azurerm_key_vault.app.vault_uri}secrets/github-private-key"
    identity            = azurerm_user_assigned_identity.app.id
  }

  template {
    min_replicas = 1
    max_replicas = 1

    container {
      name   = "container"
      image  = "${azurerm_container_registry.acr.login_server}/discordbot:${var.commit_sha256}"
      cpu    = 0.25
      memory = "0.5Gi"

      startup_probe {
        port = 8081
        transport = "HTTP"
        path = "/startup"
      }

      readiness_probe {
        port = 8081
        transport = "HTTP"
        path = "/readiness"
      }

      liveness_probe {
        port = 8081
        transport = "HTTP"
        path = "/liveness"
      }

      env {
        name  = "GITHUB_APP_ID"
        value = var.github-app-id
      }

      env {
        name  = "GITHUB_INSTALLATION_ID"
        value = var.github-installation-id
      }

      env {
        name        = "GITHUB_PRIVATE_KEY"
        secret_name = "github-private-key"
      }

      env {
        name  = "DISCORD_PUBLIC_KEY"
        value = var.discord-public-key
      }

      env {
        name  = "DISCORD_APP_ID"
        value = var.discord-app-id
      }

      env {
        name        = "DISCORD_TOKEN"
        secret_name = "discord-token"
      }
    }
  }
}

resource "azurerm_monitor_diagnostic_setting" "containerapp" {
  name = "containerapp"
  target_resource_id = azurerm_container_app_environment.app.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.app.id

  enabled_log {
    category = "ContainerAppHTTPLogs"
  }
}

resource "azurerm_user_assigned_identity" "app" {
  name                = "kokkimusume-discortbot-identity"
  location            = data.azurerm_resource_group.app.location
  resource_group_name = data.azurerm_resource_group.app.name
}

resource "azurerm_role_assignment" "acr_pull" {
  scope                = azurerm_container_registry.acr.id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_user_assigned_identity.app.principal_id
}

resource "azurerm_role_assignment" "key_vault" {
  scope                = azurerm_key_vault.app.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.app.principal_id
}
