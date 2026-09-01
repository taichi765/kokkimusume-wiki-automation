
data "azurerm_resource_group" "acr" {
  name = "kokkimusume-acr"
}

resource "azurerm_container_registry" "acr" {
  name                = "kokkimusumeDiscordbotRegistry"
  resource_group_name = data.azurerm_resource_group.acr.name
  location            = data.azurerm_resource_group.acr.location
  sku                 = "Standard"
  admin_enabled       = false
}

module "discord_bot" {
  source = "./discord_bot"

  github-installation-id = var.github-installation-id
  github-app-id          = var.github-app-id
  discord-app-id         = var.discord-app-id
  discord-public-key     = var.discord-public-key
  commit_sha256          = var.commit_sha256

  acr_group_name = data.azurerm_resource_group.acr.name
  acr_name       = azurerm_container_registry.acr.name
}

module "deletion_detector" {
  source = "./deletion_detector"

  acr_id           = azurerm_container_registry.acr.id
  acr_login_server = azurerm_container_registry.acr.login_server
  commit_hash      = var.commit_sha256
}
