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

variable "location" {
  type        = string
  description = "Location of the Bucket"
}

variable "bucket_name" {
  type        = string
  description = "Name of the bucket"
}