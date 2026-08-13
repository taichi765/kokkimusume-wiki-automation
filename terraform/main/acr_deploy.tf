data "terraform_remote_state" "boot" {
  backend = "azurerm"

  config = {
    use_azuread_auth     = true
    storage_account_name = "kdbaccount"
    container_name       = "tfstate-storage-container"
    key                  = "boot.tfstate"
  }
}

resource "azuread_application" "acr-deploy" {
  display_name = "kokkimusume-discordbot-acr-deploy"
  owners       = [data.terraform_remote_state.boot.outputs.service_principal_id]
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
