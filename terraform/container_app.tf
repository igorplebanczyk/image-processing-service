resource "azurerm_container_app_environment" "container_app_environment" {
  name = var.container_app_environment_name
  location = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_container_app" "container_app" {
  name = var.container_app_name
  resource_group_name = azurerm_resource_group.rg.name
  container_app_environment_id = azurerm_container_app_environment.container_app_environment.id
  revision_mode = var.container_app_revision_mode

  template {
    container {
      name   = var.container_app_container_name
      image  = var.container_app_container_image
      cpu    = var.container_app_container_cpu
      memory = var.container_app_container_memory
    }
  }
}