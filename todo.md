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
- [ ] Create and store image previews
  - [ ] Add a preview service for generating previews
  - [ ] Update blob names to include a preview suffix
  - [ ] Store previews in the blob storage
  - [ ] Add preview endpoints: both singular and batch
- [X] Add a mail service
  - [X] Add an admin broadcast system
- [ ] Improve security and QoL
  - [ ] Add a forgotten password system
  - [ ] Add an email verification system
  - [ ] Add 2FA
- [ ] Revise the project structure
  - [ ] Collect more logs
  - [ ] Set up a wrapper for DB queries to collect logs and metrics
  - [ ] Streamline the API, verify API security and semantics
  - [ ] Verify and improve authentication and authorization
  - [ ] Attempt to streamline the main package
  - [ ] Standardize naming conventions
  - [ ] Split large files into smaller ones
  - [ ] Extract common functionality into helper functions
  - [ ] Unit test the helper functions
  - [ ] Add comments where necessary
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

# Plan


## API

| **Category** | **Method** | **Endpoint**             | **Description**                                       |
|--------------|------------|--------------------------|-------------------------------------------------------|
| **Meta**     | ANY        | `/health`                | Check the health of the service.                      |
|              | GET        | `/metrics`               | Retrieve service metrics (private).                   |
| **Auth**     | POST       | `/auth/login/one`        | Provide credentials, get OTP.                         |
|              | POST       | `/auth/login/two`        | Provide OTP, get JWT and refresh token.               |
|              | DELETE     | `/auth/logout`           | Invalidate refresh token.                             |
|              | POST       | `/auth/refresh`          | Provide refresh token, get new JWT and refresh token. |
| **Users**    | POST       | `/users`                 | Register a new user.                                  |
|              | GET        | `/users`                 | Get details of the authenticated user.                |
|              | PUT        | `/users`                 | Update details of the authenticated user.             |
|              | DELETE     | `/users`                 | Delete the authenticated user.                        |
|              | POST       | `/users/verify`          | Provide email, get OTP for verification.              |
|              | PATCH      | `/users/verify`          | Provide OTP, verify email.                            |
|              | POST       | `/users/forgot-password` | Provide email, get OTP for password reset.            |
|              | PATCH      | `/users/forgot-password` | Provide OTP, set a new password.                      |
| **Images**   | POST       | `/images`                | Upload a new image.                                   |
|              | GET        | `/images`                | Get details of a specific image.                      |
|              | GET        | `/images/all`            | Get all images of the authenticated user.             |
|              | PUT        | `/images`                | Update image details.                                 |
|              | PATCH      | `/images`                | Apply transformations to an image.                    |
|              | DELETE     | `/images`                | Delete a specific image.                              |
| **Admin**    | GET        | `/admin/verify`          | Verify if the authenticated user is an admin.         |
|              | POST       | `/admin/broadcast`       | Send a broadcast message.                             |
|              | DELETE     | `/admin/auth/{id}`       | Logout a specific user by ID.                         |
|              | GET        | `/admin/users/{id}`      | Get details of a specific user.                       |
|              | GET        | `/admin/users`           | Get details of all users.                             |
|              | PATCH      | `/admin/users/{id}`      | Update the role of a specific user.                   |
|              | DELETE     | `/admin/users/{id}`      | Delete a specific user.                               |
|              | GET        | `/admin/images`          | Get all images (admin view).                          |
|              | DELETE     | `/admin/images`          | Delete an image (admin operation).                    |

## Structure

### Common

Cache:
- Types: Service

Database:
- Types: Service
- Functions: Connect, CloseConnection
- Subservices:
  - Tx:
    - Types: Provider
    - Functions: Transact
  - Worker:

Emails:
- Types: Service
- SendText, SendHTML, send

Errors:
- Types: Error (Type)

Logs:
- Functions: init

Metrics:
- Functions: init, Handler

Server:
- Types: Service
- Functions: Start, Stop, setup, health
- Subservices:
  - Respond:
    - Functions: WithError, WithJSON, WithNoContent, applyCommonHeaders
  - Telemetry:
    - Types: responseRecorder
    - Functions: Middleware, WriteHeader

Storage:
- Types: Service
- Subservices:
  - Worker:

### Auth

Domain:
- Types: User (Role), RefreshToken, OTP
- Repositories:
  - UserDBRepository: GetUserByUsername, GetUserRoleByID
  - RefreshTokenDBRepository: CreateRefreshToken, GetRefreshTokensByUserID, RevokeRefreshTokenByToken, RevokeRefreshTokensByUserID
  - OTPCacheRepository: AddOTP, GetOTP, DeleteOTP

Application:
- Types: AuthService
- Functions: LoginOne, LoginTwo, Logout, Refresh, generateAccessToken, verifyAccessToken, generateAccessToken, verifyRefreshToken, generateOTP, verifyOTP

Interfaces:
- Types: AuthAPI
- Functions: LoginOne, LoginTwo, Logout, Refresh, UserMiddleware, AdminMiddleware, storeTokenInCookie, extractTokenFromCookie

Infrastructure:
- Repositories:
  - UserDBRepository: GetUserByUsername, GetUserRoleByID
  - RefreshTokenDBRepository: CreateRefreshToken, GetRefreshTokensByUserID, RevokeRefreshTokenByToken, RevokeRefreshTokensByUserID
  - OTPCacheRepository: AddOTP, GetOTP, DeleteOTP


### Users

Domain:
- Types: User (Role), OTP (Type)
- Functions: ValidateUsername, ValidateEmail, ValidatePassword, DetermineUserDetailsToUpdate
- Repositories:
  - UserDBRepository: CreateUser, GetUserByID, GetAllUsers, UpdateUserDetails, UpdateUserRole, DeleteUser
  - OTPCacheRepository: AddOTP, GetOTP, DeleteOTP

Application:
- Types: UserService
- Functions: Register, Get, Update, Delete, SendUserVerificationCode, VerifyUser, SendForgotPasswordCode, ResetPassword, AdminGetAll, AdminUpdateRole, hashPassword, verifyPassword, generateOTP, verifyOTP

Interfaces:
- Types: UserAPI
- Functions: Register, Get, Update, Delete, SendUserVerificationCode, VerifyUser, SendForgotPasswordCode, ResetPassword, AdminGet, AdminGetAll, AdminUpdateRole, AdminDelete

Infrastructure:
- Repositories:
  - UserDBRepository: CreateUser, GetUserByID, GetAllUsers, UpdateUserDetails, UpdateUserRole, DeleteUser
  - OTPCacheRepository: AddOTP, GetOTP, DeleteOTP


### Images

Domain:
- Types: Image, Transformation (TransformationType)
- Functions: ValidateName, ValidateRawImage, Validate Transformation, CreateImageName, CreatePreviewName
- Repositories:
  - ImageDBRepository: CreateImage, GetImageByUserIDandName, GetImagesByUserID, GetAllImages, UpdateImage, DeleteImageByUserIDandName, DeleteImageByID
  - StorageRepository: UploadImage, DownloadImage, DeleteImage
  - CacheRepository: AddImage, GetImage, DeleteImage

Application:
- Types: ImageService
- Functions: Upload, Get, GetAll, Update, Transform, Delete, AdminGetAll, AdminDelete
- Subservices:
  - Transformations:

Interfaces:
- Types: ImageAPI
- Functions: Upload, Get, GetAll, Update, Transform, Delete, AdminGetAll, AdminDelete

Infrastructure:
- Repositories:
  - ImageDBRepository: CreateImage, GetImageByUserIDandName, GetImagesByUserID, GetAllImages, UpdateImage, DeleteImageByUserIDandName, DeleteImageByID
  - StorageRepository: UploadImage, DownloadImage, DeleteImage
  - CacheRepository: AddImage, GetImage, DeleteImage
