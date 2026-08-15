terraform {
  backend "azurerm" {
    use_azuread_auth     = true
    storage_account_name = "kdbaccount"
    container_name       = "tfstate-storage-container"
    key                  = "terraform.tfstate"
  }
}

provider "azurerm" {
  subscription_id = "ad53c5b5-a699-4d4e-8125-0619976eb022"
  features {
    key_vault {
      purge_soft_deleted_secrets_on_destroy = true
      recover_soft_deleted_secrets = true
    }
  }
}

provider "azuread" {

}
