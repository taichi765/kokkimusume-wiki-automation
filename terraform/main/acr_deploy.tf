resource "azuread_application" "acr-deploy" {
  display_name = "kokkimusume-discordbot-acr-deploy"
  owners       = ["f09650f5-9059-46d8-910f-619a2738737a"]
}

resource "azuread_service_principal" "acr-deploy" {
  client_id = azuread_application.acr-deploy.client_id
}

resource "azuread_application_federated_identity_credential" "acr-deploy" {
  application_id = azuread_application.acr-deploy.id
  display_name   = "kokkimusume-discordbot-acr-deploy-cred"
  description    = "Github Actions OIDC federation"
  audiences      = ["api://AzureADTokenExchange"]
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:taichi765@190380265/kokkimusume-wiki-automation@1328476148:ref:refs/heads/master"
}

resource "azurerm_role_assignment" "acr_push" {
  scope                = azurerm_container_registry.acr.id
  role_definition_name = "AcrPush"
  principal_id         = azuread_service_principal.acr-deploy.object_id
}
