provider "google" {
  credentials = file(var.credentials_file)
  project     = var.gcp_project
  region      = var.gcp_region
  zone        = var.gcp_az
}

####################################################
###         Create the Storage Bucket             ##
####################################################
#module "bucket" {
#  source   = "./bucket"
#  location = "NORTHAMERICA-NORTHEAST1"
#}
#
####################################################
###           Create the GKE Cluster              ##
####################################################
#module "gke_cluster" {
#  source             = "./gke"
#  initial_node_count = 2
#  machine_type       = "e2-micro"
#  name               = "image-repo-backend-cluster"
#  project            = var.gcp_project
#  region             = var.gcp_region
#}