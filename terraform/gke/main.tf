provider "google" {
  credentials = file(var.credentials_file)
  project     = var.gcp_project
  region      = var.gcp_region
  zone        = var.gcp_az
}

###################################################
##       Create the GKE container cluster        ##
###################################################
resource "google_container_cluster" "primary" {
  name     = "${var.gcp_project}-gke"
  network  = "default"
  location = var.gcp_region

  remove_default_node_pool = true
  initial_node_count       = var.initial_node_count

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }
}

###################################################
##      Create the GKE container node pool       ##
###################################################
resource "google_container_node_pool" "primary_preemptible_nodes" {
  name       = "${var.gcp_project}-node-pool"
  cluster    = google_container_cluster.primary.name
  location   = var.gcp_region
  node_count = 1

  node_config {
    preemptible  = var.preemptible
    machine_type = var.machine_type

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}