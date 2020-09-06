###################################################
##       Create the GKE container cluster        ##
###################################################
resource "google_container_cluster" "primary" {
  name     = var.name
  location = var.region

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
  name       = "${var.name}-node-pool"
  location   = var.region
  cluster    = google_container_cluster.primary.name
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