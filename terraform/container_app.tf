

resource "azurerm_log_analytics_workspace" "kokkimusume-discordbot" {
  name                = "kokkimusume-discordbot-logs"
  location            = azurerm_resource_group.kokkimusume-discordbot.location
  resource_group_name = azurerm_resource_group.kokkimusume-discordbot.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_key_vault" "kokkimusume-discordbot" {
  name                = "discordbot-key-vault"
  location            = azurerm_resource_group.kokkimusume-discordbot.location
  resource_group_name = azurerm_resource_group.kokkimusume-discordbot.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                   = "standard"
  rbac_authorization_enabled = true
}

resource "azurerm_key_vault_secret" "discord-token" {
  name         = "discord-token"
  key_vault_id = azurerm_key_vault.kokkimusume-discordbot.id
}

resource "azurerm_key_vault_secret" "github-private-key" {
  name         = "github-private-key"
  key_vault_id = azurerm_key_vault.kokkimusume-discordbot.id
}

resource "azurerm_container_app_environment" "kokkimusume-discordbot" {
  name                       = "kokkimusume-discordbot-env"
  location                   = azurerm_resource_group.kokkimusume-discordbot.location
  resource_group_name        = azurerm_resource_group.kokkimusume-discordbot.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.kokkimusume-discordbot.id
  logs_destination           = "log-analytics"
}

resource "azurerm_container_app" "kokkimusume-discordbot" {
  name                         = "kokkimusume-discordbot-app"
  container_app_environment_id = azurerm_container_app_environment.kokkimusume-discordbot.id
  resource_group_name          = azurerm_resource_group.kokkimusume-discordbot.name
  revision_mode                = "Single"

  depends_on = [azurerm_role_assignment.acr_pull, azurerm_role_assignment.key_vault]

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.kokkimusume-discordbot.id]
  }

  registry {
    server   = azurerm_container_registry.acr.login_server
    identity = azurerm_user_assigned_identity.kokkimusume-discordbot.id
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
    key_vault_secret_id = azurerm_key_vault_secret.discord-token.id
    identity            = azurerm_user_assigned_identity.kokkimusume-discordbot.id
  }

  secret {
    name                = "github-private-key"
    key_vault_secret_id = azurerm_key_vault_secret.github-private-key.id
    identity            = azurerm_user_assigned_identity.kokkimusume-discordbot.id
  }

  template {
    min_replicas = 1
    max_replicas = 1

    container {
      name   = "kokkimusume-discordbot-container"
      image  = "${azurerm_container_registry.acr.login_server}/discordbot:${var.commit_sha256}"
      cpu    = 0.25
      memory = "0.5Gi"

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

resource "azurerm_user_assigned_identity" "kokkimusume-discordbot" {
  name                = "kokkimusume-discortbot-identity"
  location            = azurerm_resource_group.kokkimusume-discordbot.location
  resource_group_name = azurerm_resource_group.kokkimusume-discordbot.name
}

resource "azurerm_role_assignment" "acr_pull" {
  scope                = azurerm_container_registry.acr.id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_user_assigned_identity.kokkimusume-discordbot.principal_id
}

resource "azurerm_role_assignment" "key_vault" {
  scope                = azurerm_key_vault.kokkimusume-discordbot.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.kokkimusume-discordbot.principal_id
}
