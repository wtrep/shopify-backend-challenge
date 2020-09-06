###################################################
##        Create the images GCP bucket           ##
###################################################
resource "google_storage_bucket" "images" {
  name     = "shopify-backend-challenge-images-bucket"
  location = var.location
}

###################################################
##     Create the GCP bucket signing account     ##
###################################################
resource "google_service_account" "bucket_service_account" {
  account_id   = "bucket-sa"
  display_name = "Service Account"
}

###################################################
##             Create the SA policies            ##
###################################################
data "google_iam_policy" "bucket_storage_admin" {
  binding {
    role = "roles/storage.admin"
    members = [
      "serviceAccount:${google_service_account.bucket_service_account.email}"
    ]
  }
}

###################################################
##     Attach the SA policies to the bucket      ##
###################################################
resource "google_storage_bucket_iam_policy" "bucket_access_iam" {
  bucket      = google_storage_bucket.images.name
  policy_data = data.google_iam_policy.bucket_storage_admin.policy_data
}

###################################################
##  Create the GCP bucket signing account Key    ##
###################################################
resource "google_service_account_key" "bucket_sa_key" {
  service_account_id = google_service_account.bucket_service_account.name
}