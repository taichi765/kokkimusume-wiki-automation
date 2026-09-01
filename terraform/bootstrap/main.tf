data "azurerm_storage_account" "backend" {
  resource_group_name = "kokkimusume-discordbot-tfstate-resources"
  name                = "kdbaccount"
}

data "azuread_application_published_app_ids" "well_known" {}

resource "azuread_service_principal" "msgraph" {
  client_id    = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
  use_existing = true
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

resource "azurerm_role_assignment" "detector_contributor" {
  scope                = azurerm_resource_group.detector.id
  role_definition_name = "Contributor"
  principal_id         = azuread_service_principal.gh-actions-apply.object_id
}

resource "azurerm_role_assignment" "storage" {
  scope                = data.azurerm_storage_account.backend.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azuread_service_principal.gh-actions-apply.object_id
}

resource "azurerm_role_assignment" "role_admin" {
  scope                = azurerm_resource_group.app.id
  role_definition_name = "User Access Administrator"
  principal_id         = azuread_service_principal.gh-actions-apply.object_id
}

resource "azuread_app_role_assignment" "readwrite" {
  app_role_id         = azuread_service_principal.msgraph.app_role_ids["Application.ReadWrite.OwnedBy"]
  principal_object_id = azuread_service_principal.gh-actions-apply.object_id
  resource_object_id  = azuread_service_principal.msgraph.object_id
}

resource "azurerm_resource_group" "app" {
  name     = "kokkimusume-discordbot-resources"
  location = "Japan East"
}

resource "azurerm_resource_group" "detector" {
  name     = "deletion-detector"
  location = "Japan East"
}
