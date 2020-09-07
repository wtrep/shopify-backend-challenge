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

variable "initial_node_count" {
  type        = number
  description = "Number of nodes to deploy"
}

variable "machine_type" {
  type        = string
  description = "Type of the nodes to deploy"
}

variable "preemptible" {
  type        = bool
  description = "Activate preemptible node pool"
  default     = false
}