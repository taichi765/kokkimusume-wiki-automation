data "azurerm_resource_group" "app" {
  name = "kokkimusume-discordbot-resources"
}

resource "azurerm_container_registry" "acr" {
  name                = "kokkimusumeDiscordbotRegistry"
  resource_group_name = data.azurerm_resource_group.app.name
  location            = data.azurerm_resource_group.app.location
  sku                 = "Standard"
  admin_enabled       = false
}
