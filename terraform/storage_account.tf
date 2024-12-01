resource "azurerm_storage_account" "storage_account" {
  name                     = var.storage_account_name
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  account_tier             = var.storage_account_tier
  account_replication_type = var.storage_account_replication_type
}

resource "azurerm_storage_container" "blob_container" {
  name                  = var.storage_account_container_name
  storage_account_name  = azurerm_storage_account.storage_account.name
  container_access_type = var.storage_account_container_access_type
}
