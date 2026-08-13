data "azurerm_storage_account" "backend" {
  resource_group_name = "kokkimusume-discordbot-tfstate-resources"
  name                = "kdbaccount"
}

resource "azurerm_resource_group" "boot" {
  name     = "kokkimusume-discordbot-boot-resources"
  location = "Japan East"
}

resource "azuread_application" "gh-actions-apply" {
  display_name = "kokkimusume-discordbot-gh-actions-appply"
}

resource "azuread_service_principal" "gh-actions-apply" {
  client_id = azuread_application.gh-actions-apply.client_id
}

resource "azuread_application_federated_identity_credential" "gh-actions-apply" {
  application_id = azuread_application.gh-actions-apply.id
  display_name   = "kokkimusume-discordbot-gh-actions-apply-cred"
  description    = "Github Actions OIDC federation for applying terraform"
  audiences      = ["api://AzureADTokenExchange"]
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:taichi765@190380265/kokkimusume-wiki-automation@1328476148:ref:refs/heads/master"
}

resource "azurerm_role_assignment" "contributor" {
  scope                = azurerm_resource_group.app.id
  role_definition_name = "Contributor"
  principal_id         = azuread_service_principal.gh-actions-apply.object_id
}


resource "azurerm_role_assignment" "storage" {
  scope                = data.azurerm_storage_account.backend.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azuread_service_principal.gh-actions-apply.object_id
}

resource "azurerm_resource_group" "app" {
  name     = "kokkimusume-discordbot-resources"
  location = "Japan East"
}
