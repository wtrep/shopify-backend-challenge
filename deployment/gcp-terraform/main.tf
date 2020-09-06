provider "google" {
  credentials = file(var.credentials_file)
  project     = var.gcp_project
  region      = var.gcp_region
  zone        = var.gcp_az
}

###################################################
##         Create the Cloud SQL Instance         ##
###################################################
resource "google_sql_database_instance" "master" {
  name             = "shopify-repo-instance"
  database_version = "MYSQL_5_7"

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      authorized_networks {
        name  = "Whitelisted CIDR block"
        value = var.authorized_cidr
      }
    }
  }
}

###################################################
##            Create the main database           ##
###################################################
resource "google_sql_database" "db" {
  name     = "images-repo"
  instance = google_sql_database_instance.master.name
}

###################################################
##            Create the test database           ##
###################################################
resource "google_sql_database" "test-db" {
  name     = "images-repo-test"
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
##         Create the Storage Bucket             ##
###################################################
module "bucket" {
  source   = "./bucket"
  location = "NORTHAMERICA-NORTHEAST1"
}

###################################################
##           Create the GKE Cluster              ##
###################################################
module "gke_cluster" {
  source             = "./gke"
  initial_node_count = 2
  machine_type       = "e2-micro"
  name               = "image-repo-backend-cluster"
  project            = var.gcp_project
  region             = var.gcp_region
}