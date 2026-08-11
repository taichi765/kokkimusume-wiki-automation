data "azurerm_client_config" "current" {
  
}

output "azure_deploy_bot_client_id" {
  value = azuread_application.github-actions.client_id
}

output "azure_tenant_id"{
  value = data.azurerm_client_config.current.tenant_id
}

output "azure_subscription_id" {
  value = data.azurerm_client_config.current.subscription_id
}

output "azure_apply_client_id" {
  value = azuread_application.gh-actions-apply.client_id
}

output "acr-login-server" {
  value = azurerm_container_registry.acr.login_server
}