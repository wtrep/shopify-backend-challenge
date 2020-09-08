# Cloud Storage Bucket

This Terraform module creates the following GCP resources:
 * A Cloud Storage bucket
 * A GCP SA to allow write access to the bucket
 * A GCP bucket signing account Key

Each variable is described in the `variable.tf` file and an example `example.tfvars` is also provided.