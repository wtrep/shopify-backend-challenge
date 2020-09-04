provider "google" {
  credentials = file(var.credentials_file)
  project     = var.gcp_project
  region      = var.gcp_region
  zone        = var.gcp_az
}

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

resource "google_sql_database" "users" {
  name     = "images-repo"
  instance = google_sql_database_instance.master.name
}

resource "google_sql_user" "db-user" {
  name     = "backend-sa"
  instance = google_sql_database_instance.master.name
  password = var.db_password
}