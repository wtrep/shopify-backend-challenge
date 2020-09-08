# Cloud SQL Instance

This Terraform module creates the following GCP resources:
 * A Cloud SQL Instance
 * A Main Database
 * A Test Database
 * A MySQL account to access the DB
 * A GCP SA to access the CloudSQL Proxy

Each variable is described in the `variable.tf` file and an example `example.tfvars` is also provided.

I was unable to bind the `roles/cloudsql.client` role to the SA using terraform. A temporary workaround is to add the role with the console or use the CLI.