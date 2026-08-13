data "azurerm_client_config" "current" {

}

output "azure_tenant_id" {
  value = data.azurerm_client_config.current.tenant_id
}

output "azure_subscription_id" {
  value = data.azurerm_client_config.current.subscription_id
}

output "acr-login-server" {
  value = azurerm_container_registry.acr.login_server
}
