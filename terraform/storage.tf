provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "rg" {
  name = var.resource_group_name
}

resource "azurerm_storage_account" "storage_account" {
  name                     = var.storage_account_name
  resource_group_name      = data.azurerm_resource_group.rg.name
  location                 = data.azurerm_resource_group.rg.location
  account_tier             = var.account_tier
  account_replication_type = var.replication_type
}

resource "azurerm_storage_container" "blob_container" {
  name                  = var.container_name
  storage_account_name  = azurerm_storage_account.storage_account.name
  container_access_type = var.container_access_type
}
