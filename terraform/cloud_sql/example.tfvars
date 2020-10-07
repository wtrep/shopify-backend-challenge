gcp_project      = "<YOUR PROJECT NAME>"
gcp_region       = "<YOUR GCP REGION>"
gcp_az           = "<YOUR GCP ZONE>"
db_instance_name = "repo-db-instance"
db_name          = "repo-db"
db_password      = "<YOUR DB PASSWORD>"
authorized_networks = [{
  "name" : "home"
  "value" = "<YOUR PUBLIC IP>/32"
}]
sql_service_account_name = "sql-sa"