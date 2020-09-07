# Shopify Backend Challenge

This personal project is my submission for the Shopify Backend Challenge Winter 2021. It consists of two microservices that I built using Go including the Terraform HCL code to provision GCP resources and the Kubernetes
YAML configuration files. The project is separated into three repositories :
 * The present [shopify-backend-challenge](https://github.com/wtrep/shopify-backend-challenge) repository that contains the files related to the deployment of the microservices. 
 * The [shopify-backend-challenge-auth](https://github.com/wtrep/shopify-backend-challenge-auth) repository that contains the source code for the authentification microservice
 * The [shopify-backend-challenge-image](https://github.com/wtrep/shopify-backend-challenge-image) repository that contains the source code for the image microservice

Each repository contain details about its content. To get detailed information about how to query the API, you can visit the [Swagger documentation page](https://app.swaggerhub.com/apis-docs/wtrep/shopify-images-repo/1.0.0).

## Terraform
The terraform folder contains three modules that create seperatly their respective content :
````
.
├── bucket     // Create the Cloud Storage Bucket and a service account with write access
├── cloud_sql  // Create the Cloud SQL instance, a production and a test MySQL database
└── gke        // Create the Kubernetes cluster
````

## Kubernetes
````
.
├── auth-microservice-deployment.yml  // Deploy the auth microservice including the sidecar CloudSQL proxy container
├── backend-env-vars-config-map.yml   // Define commonly used environment variables
├── cluster-ip-service.yml            // Allow one ClusterIP to each microservice
├── image-microservice-deployment.yml // Deploy the image microservice including the sidecar CloudSQL proxy container
└── ingress.yml                       // Configure the NGINX ingress controller and setup the http paths
````

## Possible improvements
 * Create a seperate `google_compute_network` (was planned but drop because of issues with VPC Network Peering to allow GCP resources access without going out of the internal network)
 * Assign a private IP to the Cloud SQL instance (also drop for the same reason)
 * Deploy in another Kubernetes namespace then the default one
 * Create Kubernetes secret through Terraform instead of manually
 * Automate the installation of the community managed [NGINX ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/#gce-gke)
 * Add outputs to the terraform modules to allow better reusability