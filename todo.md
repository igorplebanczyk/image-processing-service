# TODO

- [X] Improve error handling
- [X] Improve logging
- [X] Add Terraform IaC
- [X] Add a CI pipeline
- [X] Add dependabot configuration
- [X] Add RBAC
- [X] Add admin endpoints
- [ ] Add tests
  - [X] Add domain layer unit tests
  - [X] Add application layer unit tests
    - [ ] Add unit tests for subservices of the application layer
  - [ ] Add interfaces layer unit tests
  - [ ] Add infrastructure layer integration tests
  - [ ] Add component tests
  - [ ] Add end-to-end tests
- [X] Revise the API setup
  - [X] Add an endpoint for verifying whether a user is an admin
  - [X] Add a versioning system
  - [X] Streamline the error responses
  - [X] Revise the semantics of the endpoints
  - [X] Consider using path and query parameters more
- [X] Revise the logging system
  - [X] Add more logs
  - [X] Add more context to the logs
  - [X] Add a log aggregation system
  - [X] Use proper Loki and Promtail configurations
  - [X] Consider storing log indexes in a volume/cloud storage
- [X] Add metrics
  - [X] Add a metrics service
  - [X] Collect metrics across the application
  - [X] Export metrics to Prometheus via a private endpoint
  - [X] Setup Prometheus in Docker Compose
- [X] Add a dashboard
  - [X] Add Grafana to the Docker Compose setup
  - [X] Add a dashboard for the metrics
  - [X] Add a dashboard for the logs
- [X] Revise the PostgreSQL setup
  - [X] Revise the migrations system
  - [X] Streamline and revise the schema setup
  - [X] Rethink adding proper configurations
- [x] Revise the Redis setup
  - [X] Revise the Redis capabilities and how to use them properly
  - [X] Rethink adding proper configurations
- [X] Revise the Docker setup
  - [X] Streamline the Dockerfiles
  - [X] Revise naming conventions
  - [X] Consider using a more explicit network setup
  - [X] Consider using more configurations
- [X] Redesign the transformation service
- [X] Add a storage worker for deleting dangling images
- [X] Create and store image previews
  - [X] Add a preview service for generating previews
  - [X] Update blob names to include a preview suffix
  - [X] Store previews in the blob storage
  - [X] Add a preview endpoint
- [X] Add a mail service
  - [X] Add an admin broadcast system
- [X] Improve security and QoL
  - [X] Add a forgotten password system
  - [X] Add an email verification system
  - [X] Add 2FA
- [ ] Revise the project structure
  - [X] Update the workers
  - [X] Revise the transformation service
  - [X] Revise the transaction setup
  - [X] Review the emails setup; send them asynchronously
  - [X] Revise the OTP expiration system
  - [X] Review performance
  - [X] Consider moving query extraction to a global package
  - [X] Collect more logs
  - [X] Collect more metrics
  - [X] Streamline the API, verify API security and semantics
  - [X] Attempt to streamline the main package
  - [X] Standardize naming conventions
  - [X] Revise the project structure
  - [X] Extract constants into variables
  - [X] Extract common functionality into helper functions
  - [ ] Add comments where necessary
  - [ ] Update unit tests
- [ ] Add a reverse proxy
  - [ ] Add Traefik to the Docker Compose setup
  - [ ] Add rate limiting and other security features
  - [ ] Add a caching layer
  - [ ] Revise API security best practices
  - [ ] Add a simple way to manage HTTPS and domains
- [ ] Improve the Terraform setup
  - [ ] Add a setup file for the resource group and virtual network
  - [ ] Add a setup file for Azure Container Registry
  - [ ] Add a setup file for Azure Container Apps
  - [ ] Improve the setup file for Azure Blob Storage
  - [ ] Improve the overall setup (Makefile, variables, etc.)
- [ ] Add a deployment pipeline
  - [ ] Add a CD workflow
  - [ ] Figure out how to handle environment variables
- [ ] Add documentation
  - [ ] Add a .env template
  - [ ] Add a README.md file
  - [ ] Add a SECURITY.md file
  - [ ] Add a LICENSE

## API

| **Category** | **Method** | **Endpoint**                       | **Description**                                       |
|--------------|------------|------------------------------------|-------------------------------------------------------|
| **Meta**     | ANY        | `/health`                          | Check the health of the service.                      |
|              | GET        | `/metrics`                         | Retrieve service metrics (private).                   |
| **Auth**     | POST       | `/auth/login/one`                  | Provide credentials, get OTP.                         |
|              | POST       | `/auth/login/two`                  | Provide OTP, get JWT and refresh token.               |
|              | DELETE     | `/auth/logout`                     | Invalidate refresh token.                             |
|              | POST       | `/auth/refresh`                    | Provide refresh token, get new JWT and refresh token. |
| **Users**    | POST       | `/users`                           | Register a new user.                                  |
|              | GET        | `/users`                           | Get details of the authenticated user.                |
|              | PUT        | `/users`                           | Update details of the authenticated user.             |
|              | DELETE     | `/users`                           | Delete the authenticated user.                        |
|              | POST       | `/users/verify`                    | Provide email, get OTP for verification.              |
|              | PATCH      | `/users/verify`                    | Provide OTP, verify email.                            |
|              | POST       | `/users/reset-password`            | Provide email, get OTP for password reset.            |
|              | PATCH      | `/users/reset-password`            | Provide OTP, set a new password.                      |
| **Images**   | POST       | `/images`                          | Upload a new image.                                   |
|              | GET        | `/images`                          | Get details of a specific image.                      |
|              | GET        | `/images/all?page=int&limit=int`   | Get all images of the authenticated user.             |
|              | PUT        | `/images`                          | Update image details.                                 |
|              | PATCH      | `/images`                          | Apply transformations to an image.                    |
|              | DELETE     | `/images`                          | Delete a specific image.                              |
| **Admin**    | POST       | `/admin/broadcast`                 | Send a broadcast message.                             |
|              | GET        | `/admin/auth`                      | Verify if the authenticated user is an admin.         |
|              | DELETE     | `/admin/auth/{id}`                 | Logout a specific user by ID.                         |
|              | GET        | `/admin/users/{id}`                | Get details of a specific user.                       |
|              | GET        | `/admin/users?page=int&limit=int`  | Get details of all users.                             |
|              | PATCH      | `/admin/users/{id}?role=role`      | Update the role of a specific user.                   |
|              | DELETE     | `/admin/users/{id}`                | Delete a specific user.                               |
|              | GET        | `/admin/images?page=int&limit=int` | Get all images (admin view).                          |
|              | DELETE     | `/admin/images`                    | Delete an image (admin operation).                    |
