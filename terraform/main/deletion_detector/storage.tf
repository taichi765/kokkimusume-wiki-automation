
resource "azurerm_storage_account" "detector" {
  name                     = "kokkimusumedetector"
  resource_group_name      = data.azurerm_resource_group.detector.name
  location                 = data.azurerm_resource_group.detector.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "detector" {
  name                  = "page-list"
  storage_account_id    = azurerm_storage_account.detector.id
  container_access_type = "private"
}
