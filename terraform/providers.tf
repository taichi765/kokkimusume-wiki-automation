terraform {
  backend "azurerm" {
    use_oidc             = true
    use_azuread_auth     = true
    storage_account_name = "kdbaccount"
    container_name       = "tfstate-storage-container"
    key                  = "terraform.tfstate"
  }
}

provider "azurerm" {
  subscription_id = "ad53c5b5-a699-4d4e-8125-0619976eb022"
  features {

  }
}

provider "azuread" {

}
