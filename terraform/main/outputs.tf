data "azurerm_client_config" "current" {

}

output "azure_tenant_id" {
  value = data.azurerm_client_config.current.tenant_id
}

output "azure_subscription_id" {
  value = data.azurerm_client_config.current.subscription_id
}

output "acr_deploy_client_id" {
  value = azuread_application.acr-deploy.client_id
}

output "discord_bot_latest_image_tag" {
  # 将来的にprodとdevで分けたらここも分岐する
  value = var.discord_bot_commit_hash
}

output "deletion_detector_latest_image_tag" {
  value = var.deletion_detector_commit_hash
}
