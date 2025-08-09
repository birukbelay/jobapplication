# Job Application API

A robust Go-based REST API for managing job postings and applications, built with modern technologies and best practices.

## 🚀 Features

- **User Authentication & Authorization**: JWT-based auth with role-based access control
- **Job Management**: Create, update, delete, and search job postings
- **Application System**: Apply for jobs, track application status
- **File Upload**: Support for resume and document uploads via Cloudinary
- **Admin Panel**: Platform administration capabilities
- **Email Notifications**: SMTP-based email system
- **Database Migrations**: Automated database schema management
- **API Documentation**: Auto-generated OpenAPI/Swagger documentation

## default values

- because there might be problems with smtp email sending on free services i have made the default code "0000" for singup

so for signup use "0000" if you didnt get the email

## 🛠 Technology Stack

### Backend Framework

- **Go 1.23+**: Modern, performant language with excellent concurrency
- **Fiber v2**: Fast HTTP web framework inspired by Express.js
- **Huma v2**: OpenAPI 3.1 compliant API framework with automatic validation

### Database & Caching

- **PostgreSQL**: Primary database for persistent data storage
- **Redis**: In-memory caching and session management
- **GORM**: Go ORM with advanced features and database migration support

### Authentication & Security

- **JWT**: Stateless authentication with access and refresh tokens
- **Argon2**: Secure password hashing algorithm
- **Role-based Access Control**: Fine-grained permission system

### File Storage & Communication

- **Cloudinary**: Cloud-based image and file management
- **SMTP**: Email delivery system for notifications

### Development & Deployment

- **Docker**: Containerized development and deployment
- **GitHub Actions**: CI/CD pipeline for automated testing and deployment

## 📋 Prerequisites

- Go 1.23 or higher
- PostgreSQL 12+
- Redis 6+
- Docker & Docker Compose (optional, for containerized setup)

## 🔧 Local Development Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd goauth
```

### 2. Environment Configuration

Copy the example environment file and configure your settings:

create a .env file in the root directory and add only one variable
    `ENVIRONMENT=dev`
and create another file `.env.dev` and copy all the variables from .env.example,
this allows you to have different configurations for different environments.
you can simply change env by updating the .env file and settiong the ENVIRONMENT variable to the desired environment.

### 3. Configure Environment Variables

Edit `.env` and `.env.dev` with your specific configuration:

#### Required Environment Variables

```bash
# Server Configuration
ENVIRONMENT=dev
SERVER_PORT=8001
SERVER_HOST=localhost

# PostgreSQL Database
SQL_DB_NAME=your_database_name
SQL_USERNAME=your_db_user
SQL_PASSWORD=your_db_password
SQL_HOST=localhost
SQL_PORT=5432
SQL_DRIVER=postgres
SSL_MODE=disable

# Redis Configuration
KV_HOST=localhost
KV_PORT=6379
KV_DB=0
KV_USER=default
KV_PASSWORD=your_redis_password

# JWT Configuration
ACCESS_SECRET=your_jwt_access_secret_key_here
ACCESS_SECRET_EXPIRE_MIN=600
REFRESH_SECRET=your_jwt_refresh_secret_key_here
REFRESH_SECRET_EXPIRES_MIN=6000

# Cloudinary (File Upload)
CLOUDINARY_API_SECRET=your_cloudinary_api_secret
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_CLOUD_NAME=your_cloudinary_cloud_name
CLOUDINARY_FOLDER=your_app_folder

# SMTP Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_PWD=your_email_app_password
SMTP_USERNAME=your_email@gmail.com

# S3 Configuration (Optional)
S3_ENDPOINT=your_s3_endpoint
S3_WEB_ENDPOINT=your_s3_web_endpoint
S3_REGION=your_s3_region
S3_FORCE_PATH_STYLE=true
S3_BUCKET_NAME=your_bucket_name
S3_ACCESS_KEY_ID=your_s3_access_key
S3_SECRET_ACCESS_KEY=your_s3_secret_key
```

### 4. Database Setup

#### Option A: Local PostgreSQL & Redis

1. Install and start PostgreSQL and Redis on your system
2. Create a database with the name specified in `SQL_DB_NAME`
3. The application will automatically run migrations on startup

#### Option B: Docker Compose (Recommended)

```bash
# Start the application with all dependencies
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop the application
docker-compose down
```

### 5. Manual Setup (Without Docker)

```bash
# Install dependencies
go mod download

# Run database migrations (automatic on startup)
# Build and run the application
go build -o main .
./main
```

The API will be available at `http://localhost:8001`

## 📚 API Documentation

Once the application is running, you can access the interactive API documentation at:

- **Swagger UI**: `http://localhost:8001/docs`
- **OpenAPI Spec**: `http://localhost:8001/openapi.json`

## 🔐 API Authentication

The API uses JWT-based authentication with two types of tokens:

1. **Access Token**: Short-lived (10 hours by default) for API requests
2. **Refresh Token**: Long-lived (100 hours by default) for token renewal

### Authentication Flow

1. **Register/Login**: `POST /api/v1/auth/login` or `POST /api/v1/auth/register`
2. **Include Token**: Add `Authorization: Bearer <access_token>` header to requests
3. **Refresh Token**: Use `POST /api/v1/auth/refresh` when access token expires

## 🎯 API Endpoints Overview

### Authentication

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - User logout

### User Profile

- `GET /api/v1/profile` - Get user profile
- `PATCH /api/v1/profile` - Update user profile

### Jobs

- `GET /api/v1/jobs` - List jobs (with filtering and pagination)
- `POST /api/v1/jobs` - Create job posting (Owner only)
- `GET /api/v1/jobs/{id}` - Get specific job
- `PATCH /api/v1/jobs/{id}` - Update job (Owner only)
- `DELETE /api/v1/jobs/{id}` - Delete job (Owner only)

### Applications

- `GET /api/v1/applications` - List applications (Owner/Admin only)
- `POST /api/v1/applications` - Apply for job
- `GET /api/v1/applications/{id}` - Get specific application
- `PATCH /api/v1/applications/{id}` - Update application
- `DELETE /api/v1/applications/{id}` - Withdraw application
- `GET /api/v1/my-applications` - Get current user's applications

### File Upload

- `POST /api/v1/upload` - Upload files (resumes, documents)

### Admin (Platform Admin only)

- `GET /api/v1/admin` - List  users
- `GET /api/v1/admin/{id}` - Get  user
- `PATCH /api/v1/admin/{id}` - Update  user

## 👥 User Roles

- **APPLICANT**: Can apply for jobs, manage own applications
- **COMPANY**: Can create and manage job postings, view applications
- **PLATFORM_ADMIN**: Full system access, user management

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## 🚀 Deployment

### Docker Deployment

```bash
# Build and deploy
docker-compose up -d --build

# Production deployment with environment-specific configs
docker-compose -f docker-compose.prod.yml up -d
```

### Manual Deployment

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -o main .

# Run with production environment
ENVIRONMENT=prod ./main
```

## 🔧 Configuration

The application uses a hierarchical configuration system:

1. Environment variables
2. `.env` file
3. `.env.<ENVIRONMENT>` file (e.g., `.env.dev`, `.env.prod`)

## 📝 Development Notes

### Project Structure

```
├── src/
│   ├── app/                 # Application routes and handlers
│   │   ├── account/         # Authentication and profile
│   │   ├── admin/           # Admin management
│   │   ├── job/             # Job and application management
│   │   └── general/         # General utilities (upload, etc.)
│   ├── models/              # Data models and database schemas
│   ├── providers/           # Service providers and dependencies
│   └── server/              # Server configuration and routing
├── public/                  # Static files and templates
├── docker-compose.yml       # Development environment
└── main.go                  # Application entry point
```

---

## 🏗 Technical Implementation & Design Decisions

This section documents the key techniques and architectural decisionsure performance, security, and scalability.

### 🚀 Performance & Scalability

#### Database Optimization

- **Connection Pooling**: GORM with PostgreSQL connection pooling to efficiently manage database connections
- **Indexed Queries**: Strategic database indexing on frequently queried fields (user emails, job IDs, application relationships)
- **Query Optimization**:
  - Selective field loading with `SelectedFields` paramereduce data transfer
  - Pagination with offset-based queries to handle large datasets
  - Lazy loading of relationships to prevent N+1 query problems

#### Caching Strategy

- **Redis Integration**: In-memory caching for:
  - JWT token blacklisting for logout functionality
  - Session management and user authentication state
  - Frequently accessed data to reduce database load
  - Rate limiting counters and temporary data storage

#### API Performance

- **Fiber Framework**: High-performance HTTP framework with minimal memory footprint
- **JSON Optimization**: Efficient JSON serialization with selective field marshaling
- **Generic Controllers**: generic CRUD operations to reduce code duplication and improve maintainability
- **Structured logs**: Efficient logging with structured formats for better performance monitoring

#### Scalability Design

- **Stateless Architecture**: JWT-based authentication eli server-side session storage
- **Microservice-Ready**: Modular structure with clear separatiocerns
- **Clean archtecture**: because of use of clean archtecture you can swap one module for another
- **Environment-Specific Configuration**: Easy scaling across different environments

### 🛡 Security Implementatio Authentication & Authorization

- **JWT Security**:
  - Separate access and refresh tokens with different expiration times
  - Secure tokewith configurable secrets
  - Token blacklisting on logout to prevent replay attacks
- **Password Security
  - Argon2 hashing algorithm (industry standard for password security)
  - Configurable hash parameters for future-proofing
  - Password complexity requirements through validation

#### Role-Based Access Control **Granular Permissions**: Fine-grained access control with operation-specific role requirements

- **Middleware Authorization**: Centralized authorization middleware for consistent security enforcement
- **Resource Ownership**: Users can only access and modify their own resources

#### Data Protection

- **Input Validation**:
  - Automatic request validation using Huma v2 schema validation
injection prevention through GORM's prepared statements
  - XSS protection through proper input sanitization
- **Sensitive Data Handling**:
  - Password fields excluded from JSON responses (`json:"-"`)
  - Email verification system for account security
  - Secure file upload with type validatio# API Security
- **HTTPS Enforcement**: Production-ready TLS configuration
- **CORS Configuration**: Proper cross-origin resource sharing setup
- **Request Size Limits**: Protection against large payload attacks
- **Secure Headers**: Implementation of security headers for web protection

### 🔒 Abuse Prevention




#### Request Validation & Filtering

- **Input Sanitization**: Comprehensive input validation and sanitization
- **Request Size Limits**: Maximum payload size enforcement
- **Content Type Validation**: Strict content type ecking
- **Malicious Pattern Detection**: Basic pattern matching for common attack vectors

#### Resource Protection

- **Database Query Limits**: Maximum result set sizes to prevent resource exhaustion
- **File Upload Restrictions**:
  - File size limits for uploads
  - File type validation and sanitizati

### 🔧 Monitoring & Observability

#### Logging Strategy

- **Structured Logging**: JSON-formatted logs for easy parsing and analysis
- **Request Tracing**: Unique request IDs for tracking requests across services
- **Error Tracking**: Comprehensive errorng with stack traces





### 🏛 ArchitecturaPatterns

#### Clean Architecture

- **Sepation of Concerns**: Clear boundaries between business logic, data access, and presentation
- **Dependency Injection**: Using provider pattern for better dependency management, and testing
- **Generic Patterns**: Reusable generic controllers and services to reduce code duplicaion

#### Error Handling

- **Centralized Error Handling: Consistent error response format across all endpoints
- **Graceful Degradation**: Fallbacmechanisms for external service failures
- **User-Friendly Messages**: Clear, actionable error ges for API consumers

#### Configuration Man **Environment-Based Config**: Hierarchical configuration system for different deployment environmen

- **Secret Management**: Secure handling of sensitive configuration data
- **Feature Flags**: Ready for feature toggle implementation




### 🔄 Future Scalability Considerations

- **Horizontal Scaling**: Stateless design ready for load balancer distribution
- **Database Sharding**: Architecture prepared for database partitioning
- **Caching Layers**: Multiple caching levels (Redis, CDN, application-level)


