# Kubernetes

This repository contains all the required files for the microservices deployment except:
1. You will need to manually install NGINX Community Ingress Controller (Details in the [installation guide](https://kubernetes.github.io/ingress-nginx/deploy/))
2. You will need to add manually using `kubectl` the following secrets (since base64 isn't encryption) :
   * sql-sa-secret: content of the `key.json` for the CloudSQL service account obtain by this [gcloud command](https://cloud.google.com/sql/docs/mysql/connect-kubernetes-engine#service_account_key_file).
   * cloudsql-db-password: key value where the value is the password of the CloudSQL Database
   * jwt-key: key value where the value is the JWT signing key shared by the two microservices
   * gcp-image-sa: content of the `service_account.json` for the Cloud Storage service account.

You can apply each file using:
````
kubectl apply -f <PATH OF THE FILE>
````
