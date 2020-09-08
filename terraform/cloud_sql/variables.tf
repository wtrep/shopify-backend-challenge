variable "credentials_file" {
  description = "Relative path to the GCP credentials file."
  type        = string
}

variable "gcp_project" {
  description = "GCP Project in which to deploy resources"
  type        = string
}

variable "gcp_region" {
  description = "Region in which the resources are deployed"
  type        = string
}

variable "gcp_az" {
  description = "AZ in which the resources are deployed"
  type        = string
}

variable "db_instance_name" {
  description = "Name of the db instance"
  type        = string
}

variable "db_name" {
  description = "Name of the db"
  type        = string
}

variable "db_password" {
  description = "Password for the DB service account"
  type        = string
}

variable "database_tier" {
  description = "Tier of the Database"
  type        = string
  default     = "db-f1-micro"
}

variable "database_version" {
  description = "Version of the Database"
  type        = string
  default     = "MYSQL_5_7"
}

variable "authorized_networks" {
  description = "Key-value pair of each authorized CIDR"
  type        = list(map(string))
  default     = null
}

variable "sql_service_account_name" {
  description = "Name of the Cloud SQL Service account"
  type        = string
}