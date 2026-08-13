output "service_principal_id" {
  value = azuread_service_principal.gh-actions-apply.object_id
}
