
resource "azurerm_user_assigned_identity" "detector" {
  name                = "kokkimusume-deletion-detector-identity"
  location            = data.azurerm_resource_group.detector.location
  resource_group_name = data.azurerm_resource_group.detector.name
}

resource "azurerm_role_assignment" "detector_acr_pull" {
  scope                = var.acr_id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_user_assigned_identity.detector.principal_id
}

resource "azurerm_role_assignment" "detector_key_vault" {
  scope                = azurerm_key_vault.detector.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.detector.principal_id
}

resource "azurerm_role_assignment" "detector_blob_storage" {
  scope                = azurerm_storage_container.detector.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.detector.principal_id
}
