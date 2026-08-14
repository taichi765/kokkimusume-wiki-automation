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

output "acr_deploy_client_id" {
  value = azuread_application.acr-deploy.client_id
}

output "latest_image_tag" {
  # 将来的にprodとdevで分けたらここも分岐する
  value = var.commit_sha256
}
