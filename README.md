# Golang REST API Boilerplate

A production-ready REST API boilerplate built with Go, Gin, PostgreSQL, and clean architecture principles.

## Features

- **Clean Architecture** - Separation of concerns with Handler → Service → Repository layers
- **JWT Authentication** - Secure authentication with access and refresh tokens
- **Role-Based Access Control (RBAC)** - Flexible permission system with multiple roles
- **UUID Primary Keys** - All IDs use UUIDs instead of auto-increment integers
- **Type-Safe Database Queries** - SQLC generates Go code from SQL
- **Database Migrations** - Version-controlled schema changes with golang-migrate
- **Docker Support** - Complete Docker development environment
- **Comprehensive Middleware** - CORS, Auth, RBAC, Logging, Panic Recovery
- **Standardized JSON Responses** - Consistent API response format
- **Input Validation** - Request validation with go-playground/validator
- **Prepared for Scaling** - Redis and S3/R2 integration ready

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- **Database**: PostgreSQL 15 - Reliable relational database
- **Query Builder**: [SQLC](https://sqlc.dev/) - Type-safe SQL queries
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) - Database schema versioning
- **Authentication**: JWT (JSON Web Tokens) - Stateless authentication
- **Caching**: Redis 7 (prepared for implementation)
- **Storage**: AWS S3 / Cloudflare R2 (prepared for implementation)

## Project Structure

```
xnoia-go-boilerplate/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── handler/                 # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── admin_handler.go
│   │   ├── health_handler.go
│   │   └── router.go
│   ├── service/                 # Business logic layer
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   └── admin_service.go
│   ├── repository/              # Data access layer
│   │   ├── repository.go
│   │   └── user_repository.go
│   ├── middleware/              # HTTP middleware
│   │   ├── cors.go
│   │   ├── auth.go
│   │   ├── rbac.go
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── model/                   # Domain models & DTOs
│   │   ├── user.go
│   │   ├── request.go
│   │   ├── response.go
│   │   └── error.go
│   └── utils/                   # Utility functions
│       ├── response.go
│       ├── jwt.go
│       ├── password.go
│       └── validator.go
├── db/
│   ├── migration/               # SQL migration files
│   │   ├── 000001_init_schema.up.sql
│   │   └── 000001_init_schema.down.sql
│   ├── query/                   # SQLC query definitions
│   │   ├── users.sql
│   │   ├── user_details.sql
│   │   ├── roles.sql
│   │   └── user_roles.sql
│   └── sqlc/                    # Generated SQLC code
├── pkg/
│   ├── database/                # Database connection
│   │   └── postgres.go
│   └── cache/                   # Redis client (prepared)
│       └── redis.go
├── docker-compose.yml           # Docker services
├── Dockerfile                   # Application container
├── Makefile                     # Development commands
├── sqlc.yaml                    # SQLC configuration
├── .env.example                 # Environment template
└── README.md
```

## Quick Start

### Prerequisites

- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Docker & Docker Compose** - [Download](https://www.docker.com/products/docker-desktop)
- **golang-migrate** - [Installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- **SQLC** - [Installation](https://docs.sqlc.dev/en/latest/overview/install.html)

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd xnoia-go-boilerplate
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Docker services**
   ```bash
   make docker-up
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Generate SQLC code**
   ```bash
   make sqlc
   ```

6. **Download Go dependencies**
   ```bash
   go mod tidy
   ```

7. **Start the server**
   ```bash
   make server
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Public Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login with email/username |
| POST | `/api/v1/auth/refresh` | Refresh access token |

### Protected Routes (JWT Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/me` | Get current user info |
| GET | `/api/v1/profile` | Get user profile |
| PUT | `/api/v1/profile` | Update user profile |
| POST | `/api/v1/profile/avatar` | Upload avatar image |

### Admin Routes (Admin Role Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/users` | List all users (paginated) |
| PUT | `/api/v1/admin/users/:id/roles` | Assign roles to user (`:id` is UUID) |
| DELETE | `/api/v1/admin/users/:id` | Deactivate user (`:id` is UUID) |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Check API and database health |

See [API.md](API.md) for detailed API documentation with request/response examples.

## Development

### Makefile Commands

```bash
# Database migrations
make migrate-up          # Run all migrations
make migrate-down        # Rollback last migration
make migrate-create name=migration_name  # Create new migration

# Code generation
make sqlc                # Generate SQLC code from SQL queries

# Development
make server              # Run the development server
make test                # Run tests

# Docker
make docker-up           # Start PostgreSQL + Redis
make docker-down         # Stop containers
make docker-build        # Build application Docker image

# Cleanup
make clean               # Remove generated files
```

### Running with Docker

**Development with hot reload:**
```bash
# Start services
make docker-up

# In another terminal, run the app
make server
```

**Production Docker build:**
```bash
make docker-build
docker run -p 8080:8080 --env-file .env xnoia-go-boilerplate:latest
```

### Database Migrations

**Create a new migration:**
```bash
make migrate-create name=add_user_preferences
```

**Apply migrations:**
```bash
make migrate-up
```

**Rollback migrations:**
```bash
make migrate-down
```

## Authentication Flow

1. **Register**: `POST /api/v1/auth/register`
   - Creates user account with hashed password
   - Assigns default "user" role
   - Returns access token and refresh token

2. **Login**: `POST /api/v1/auth/login`
   - Accepts email OR username + password
   - Validates credentials
   - Returns access token (24h) and refresh token (7 days)

3. **Protected Requests**: Add JWT to Authorization header
   ```
   Authorization: Bearer <access_token>
   ```

4. **Refresh Token**: `POST /api/v1/auth/refresh`
   - When access token expires, use refresh token
   - Returns new access token and refresh token

## Role-Based Access Control

### Default Roles

The system includes three default roles with fixed UUIDs:

- **user** (`550e8400-e29b-41d4-a716-446655440001`) - Standard user (default for new registrations)
- **admin** (`550e8400-e29b-41d4-a716-446655440002`) - Full administrative access
- **moderator** (`550e8400-e29b-41d4-a716-446655440003`) - Moderator with extended permissions

### Assigning Roles

Only admins can assign roles:
```bash
PUT /api/v1/admin/users/:id/roles
{
  "role_ids": ["550e8400-e29b-41d4-a716-446655440002"]
}
```

### Adding Custom Roles

Roles are stored in the `roles` table with JSONB permissions field for flexibility.

## Database Schema

### Users Table
- **Primary Key**: UUID (generated automatically)
- Stores authentication credentials
- Email and username (both unique)
- Password hash (bcrypt)
- Active status and email verification

### User Details Table
- **Primary Key**: UUID (generated automatically)
- **Foreign Key**: User UUID
- Profile information (full name, phone, address, etc.)
- Avatar URL
- Date of birth
- Bio

### Roles Table
- **Primary Key**: UUID (fixed for default roles, generated for custom roles)
- Role definitions
- JSONB permissions field

### User Roles Table
- **Primary Key**: UUID (generated automatically)
- Many-to-many relationship using UUIDs
- Tracks who assigned the role and when

**Note**: All IDs use UUIDs instead of auto-increment integers for better security and distributed system compatibility.

## Security Features

- **Password Hashing** - bcrypt with cost factor 12
- **JWT Tokens** - Signed with HS256 algorithm
- **SQL Injection Protection** - SQLC generates parameterized queries
- **Input Validation** - go-playground/validator
- **CORS** - Configurable cross-origin resource sharing
- **Panic Recovery** - Graceful error handling
- **Rate Limiting** - Ready for Redis implementation

## Testing

```bash
# Run all tests
make test

# Run with coverage
go test -v -cover ./...

# Run specific package
go test -v ./internal/service
```

## Deployment

### Environment Variables

Ensure all required environment variables are set:
- `JWT_SECRET` - Strong secret key for JWT signing
- `DB_*` - Database connection details
- `SERVER_PORT` - API server port

### Production Checklist

- [ ] Set `GIN_MODE=release`
- [ ] Use strong `JWT_SECRET`
- [ ] Enable HTTPS
- [ ] Configure CORS for your domain
- [ ] Set up database backups
- [ ] Configure logging aggregation
- [ ] Enable rate limiting
- [ ] Set up monitoring (health checks)

## Future Enhancements

### Redis Integration
- Session storage
- Rate limiting
- Cache frequently accessed data

### S3/Cloudflare R2 Integration
- Avatar uploads
- File attachments
- Presigned URLs

### Additional Features
- Email verification
- Password reset flow
- Two-factor authentication (2FA)
- OAuth2 integration
- WebSocket support
- Background job processing

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions, please open an issue on GitHub.

---

Built with ❤️ using Go
