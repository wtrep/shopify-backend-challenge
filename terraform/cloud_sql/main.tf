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
  name             = var.db_name
  database_version = "MYSQL_5_7"

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled    = true
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
