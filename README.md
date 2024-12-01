# Image Processing Service
This is an API for a simple image processing service. It allows users to upload images, transform them, and store them in the cloud. Written in Go.

## Features
* User registration and basic account management
* Two-factor authentication using JWTs and TOTPs
* Image uploading and downloading to and from Azure Blob Storage
* Image transformation using the [imaging](https://github.com/disintegration/imaging) package
* Image preview generation
* Email verification, password reset and 2FA using TOTPs
* Observability using Loki, Prometheus, and Grafana
* Reverse proxy using Traefik
* PostgreSQL and Redis for data storage
* Terraform IaC for Azure resources, including Azure Container Apps
* Docker Compose for container management
* CI/CD using GitHub Actions

## Installation
1. Clone the repository
2. Configure the environment variables in the `.env` file; see the `.env.example` file for reference
3. Make sure Docker is installed
4. Run `docker-compose up --build` to start the service
5. The API can now be accessed at `localhost:80`

## Usage
The API is documented [here](https://www.postman.com/science-meteorologist-78724576/image-processing-service/api/5b6b7130-04d2-4788-9c26-babd4d9dfbd4) using Postman.

## Dependencies

### Go packages
* [azblob](https://github.com/Azure/azure-sdk-for-go)
* [imaging](https://github.com/disintegration/imaging)
* [redis](https://github.com/go-redis/redis)
* [jwt](https://github.com/golang-jwt/jwt)
* [uuid](https://github.com/google/uuid)
* [pq](https://github.com/lib/pq)
* [otp](https://github.com/pquerna/otp)
* [prometheus](https://github.com/prometheus/client_golang)
* [go-mail](https://github.com/wneessen/go-mail)
* [x/crypto](https://golang.org/x/crypto)

### External services
* [PostgreSQL](https://www.postgresql.org/)
* [Redis](https://redis.io/)
* [Loki](https://grafana.com/oss/loki/)
* [Prometheus](https://prometheus.io/)
* [Grafana](https://grafana.com/)
* [Traefik](https://traefik.io/)
* [Azure Blob Storage](https://azure.microsoft.com/en-us/products/storage/blobs)

## License
This project is licensed under the MIT License - see [LICENSE](https://github.com/igorplebanczyk/image-processing-service/blob/main/LICENSE).

## Notes
* This project is primarily a learning exercise and is not intended for production use, though I did my best to make it at least somewhat viable for small-scale deployment.
* The project is my [boot.dev Capstone Project](https://www.boot.dev/courses/build-capstone-project) and my solution to the [roadmap.sh project](https://roadmap.sh/projects/image-processing-service).
