
resource "azurerm_storage_account" "detector" {
  name                     = "detector-storage"
  resource_group_name      = azurerm_resource_group.detector.name
  location                 = azurerm_resource_group.detector.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "detector" {
  name                  = "page-list"
  storage_account_id    = azurerm_storage_account.detector.id
  container_access_type = "private"
}
