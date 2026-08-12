
resource "azuread_application" "github-actions" {
  display_name = "kokkimusume-discordbot-github-actions"
  owners       = [azuread_application.gh-actions-apply.object_id]
}

resource "azuread_service_principal" "github-actions" {
  client_id = azuread_application.github-actions.client_id
}

resource "azuread_application_federated_identity_credential" "github-actions" {
  application_id = azuread_application.github-actions.id
  display_name   = "kokkimusume-discordbot-github-actions-cred"
  description    = "Github Actions OIDC federation"
  audiences      = ["api://AzureADTokenExchange"]
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:taichi765@190380265/kokkimusume-wiki-automation@1328476148:ref:refs/heads/master"
}

resource "azurerm_role_assignment" "acr_push" {
  scope                = azurerm_container_registry.acr.id
  role_definition_name = "AcrPush"
  principal_id         = azuread_service_principal.github-actions.object_id
}
