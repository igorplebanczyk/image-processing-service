variable "resource_group_name" {
  description = "The name of the resource group where the storage account will be created."
  type        = string
}

variable "storage_account_name" {
  description = "The name of the storage account. Must be globally unique."
  type        = string
}

variable "account_tier" {
  description = "The performance tier of the storage account. Valid options are Standard or Premium."
  type        = string
  default     = "Standard"
}

variable "replication_type" {
  description = "The replication strategy for the storage account. Valid options are LRS, GRS, RAGRS, ZRS."
  type        = string
  default     = "LRS"
}

variable "container_name" {
  description = "The name of the blob container."
  type        = string
}

variable "container_access_type" {
  description = "The access level of the blob container. Valid options are private, blob, or container."
  type        = string
  default     = "private"
}
