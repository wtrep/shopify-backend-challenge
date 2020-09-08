provider "google" {
  credentials = file(var.credentials_file)
  project     = var.gcp_project
  region      = var.gcp_region
  zone        = var.gcp_az
}

###################################################
##           Fetch the default network           ##
###################################################
data "google_compute_network" "default" {
  name = "default"
}

###################################################
##         Create the Cloud SQL Instance         ##
###################################################
resource "google_sql_database_instance" "master" {
  name             = var.db_instance_name
  database_version = var.database_version

  settings {
    tier = var.database_tier
    ip_configuration {
      ipv4_enabled    = true
      private_network = data.google_compute_network.default.self_link

      dynamic "authorized_networks" {
        for_each = var.authorized_networks
        content {
          name  = authorized_networks.value["name"]
          value = authorized_networks.value["value"]
        }
      }
    }
  }
}

###################################################
##            Create the main database           ##
###################################################
resource "google_sql_database" "db" {
  name     = var.db_name
  instance = google_sql_database_instance.master.name
}

###################################################
##            Create the test database           ##
###################################################
resource "google_sql_database" "test-db" {
  name     = "${var.db_name}-test"
  instance = google_sql_database_instance.master.name
}

###################################################
##   Create the image MySQL service account    ##
###################################################
resource "google_sql_user" "db-user" {
  name     = "backend-sa"
  instance = google_sql_database_instance.master.name
  password = var.db_password
}

###################################################
##         Create the GCP service account        ##
###################################################
resource "google_service_account" "sql-sa" {
  account_id   = var.sql_service_account_name
  display_name = var.sql_service_account_name
  description  = "GCP SQL Service Account"
}