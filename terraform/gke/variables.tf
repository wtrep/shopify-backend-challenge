variable "name" {
  type        = string
  description = "Name of the GKE cluster"
}

variable "project" {
  type        = string
  description = "GCP project to deploy in"
}

variable "region" {
  type        = string
  description = "GCP region to deploy in"
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