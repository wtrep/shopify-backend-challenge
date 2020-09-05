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

variable "db_password" {
  description = "Password for the DB service account"
  type        = string
}

variable "authorized_cidr" {
  description = "Authorized CIDR block that can access the MySQL DB"
  type        = string
}