resource "azurerm_resource_group" "kokkimusume-discordbot" {
  name     = "kokkimusume-discordbot-resources"
  location = "Japan East"
}

resource "azurerm_container_registry" "acr" {
  name = "kokkimusumeDiscordbotRegistry"
  resource_group_name = azurerm_resource_group.kokkimusume-discordbot.name
  location = azurerm_resource_group.kokkimusume-discordbot.location
  sku = "Standard"
  admin_enabled = false
}
