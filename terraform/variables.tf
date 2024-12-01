variable "location" {
    description = "The location/region where the services will be created."
    type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group."
  type        = string
}

# Storage Account

variable "storage_account_name" {
  description = "The name of the storage account. Must be globally unique."
  type        = string
}

variable "storage_account_tier" {
  description = "The performance tier of the storage account. Valid options are Standard or Premium."
  type        = string
}

variable "storage_account_replication_type" {
  description = "The replication strategy for the storage account. Valid options are LRS, GRS, RAGRS, ZRS."
  type        = string
}

variable "storage_account_container_name" {
  description = "The name of the blob container."
  type        = string
}

variable "storage_account_container_access_type" {
  description = "The access level of the blob container. Valid options are private, blob, or container."
  type        = string
}

# Container App

variable "container_app_environment_name" {
  description = "The name of the container app environment."
  type        = string
}

variable "container_app_name" {
  description = "The name of the container app."
  type        = string
}

variable "container_app_revision_mode" {
  description = "The revision mode of the container app."
  type        = string
}

variable "container_app_container_name" {
  description = "The name of the container."
  type        = string
}

variable "container_app_container_image" {
  description = "The image of the container."
  type        = string
}

variable "container_app_container_cpu" {
  description = "The CPU of the container."
  type        = number
}

variable "container_app_container_memory" {
  description = "The memory of the container."
  type        = number
}
