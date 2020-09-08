# Kubernetes

This repository contains all the required files for the microservices deployment except:
    * You will need to manually install NGINX Community Ingress Controller (Details in the [installation guide](https://kubernetes.github.io/ingress-nginx/deploy/))
    * You will need to add manually using `kubectl` the following secrets (since base64 isn't encryption) :
        * sql-sa-secret: content of the `service_account.json` for the CloudSQL service account.
        * DB_PASSWORD: key value where the value is the password of the CloudSQL Database
        * JWT_KEY: key value where the value is the JWT signing key shared by the two microservices
        * gcp-image-sa: content of the `service_account.json` for the Cloud Storage service account.

You can apply each file using:
````
kubectl apply -f <PATH OF THE FILE>
````