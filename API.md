# API Documentation

Base URL: `http://localhost:8080/api/v1`

All endpoints return JSON responses with the following structure:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Operation failed",
  "error": "Error details"
}
```

---

## Authentication Endpoints

### Register User

Creates a new user account.

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePass123",
  "full_name": "John Doe"
}
```

**Validation Rules:**
- `email`: Required, valid email format
- `username`: Required, 3-30 characters, alphanumeric + underscore
- `password`: Required, minimum 8 characters, must contain uppercase, lowercase, and number
- `full_name`: Optional

**Success Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "email_verified": false,
      "is_active": true,
      "roles": [
        {
          "id": 1,
          "role_name": "user",
          "description": "Standard user role"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

**Error Responses:**
- `400` - Validation error
- `409` - Email or username already exists

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "SecurePass123",
    "full_name": "John Doe"
  }'
```

---

### Login

Authenticates a user with email or username.

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**
```json
{
  "login": "user@example.com",
  "password": "SecurePass123"
}
```

**Notes:**
- `login` can be either email OR username
- Password is case-sensitive

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "email_verified": false,
      "is_active": true,
      "last_login": "2024-01-15T10:30:00Z",
      "roles": [
        {
          "id": 1,
          "role_name": "user",
          "description": "Standard user role"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

**Error Responses:**
- `401` - Invalid credentials
- `403` - User account is inactive

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "johndoe",
    "password": "SecurePass123"
  }'
```

---

### Refresh Token

Generates new access and refresh tokens.

**Endpoint:** `POST /api/v1/auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": { ... }
  }
}
```

**Error Responses:**
- `401` - Invalid or expired refresh token
- `403` - User account is inactive

---

### Get Current User

Retrieves the authenticated user's information.

**Endpoint:** `GET /api/v1/auth/me`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "email_verified": false,
    "is_active": true,
    "last_login": "2024-01-15T10:30:00Z",
    "roles": [
      {
        "id": 1,
        "role_name": "user",
        "description": "Standard user role"
      }
    ],
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- `401` - Missing or invalid token
- `404` - User not found

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## User Profile Endpoints

All profile endpoints require JWT authentication.

### Get Profile

Retrieves the user's complete profile including details.

**Endpoint:** `GET /api/v1/profile`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "email_verified": false,
      "is_active": true,
      "last_login": "2024-01-15T10:30:00Z",
      "roles": [
        {
          "id": 1,
          "role_name": "user",
          "description": "Standard user role"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z"
    },
    "details": {
      "full_name": "John Doe",
      "phone": "+1234567890",
      "address": "123 Main St, City, Country",
      "date_of_birth": "1990-01-15",
      "avatar_url": "/avatars/user-1-1234567890.jpg",
      "bio": "Software developer passionate about Go"
    }
  }
}
```

---

### Update Profile

Updates user profile details.

**Endpoint:** `PUT /api/v1/profile`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "full_name": "John Doe Updated",
  "phone": "+1234567890",
  "address": "456 New St, City, Country",
  "date_of_birth": "1990-01-15",
  "bio": "Updated bio text"
}
```

**Notes:**
- All fields are optional
- Only provided fields will be updated
- `date_of_birth` format: `YYYY-MM-DD`
- `bio` max length: 500 characters

**Success Response (200):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": null
}
```

**Error Responses:**
- `400` - Validation error
- `401` - Unauthorized
- `404` - User not found

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe Updated",
    "bio": "Software engineer"
  }'
```

---

### Upload Avatar

Uploads a user avatar image.

**Endpoint:** `POST /api/v1/profile/avatar`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

**Request Body:**
- Form field name: `avatar`
- Accepted types: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- Max size: 5MB

**Success Response (200):**
```json
{
  "success": true,
  "message": "Avatar updated successfully",
  "data": {
    "avatar_url": "/avatars/user-1-1705318200.jpg"
  }
}
```

**Error Responses:**
- `400` - Invalid file type or file too large
- `401` - Unauthorized

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/profile/avatar \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "avatar=@/path/to/image.jpg"
```

---

## Admin Endpoints

All admin endpoints require JWT authentication AND admin role.

### List Users

Retrieves a paginated list of all users.

**Endpoint:** `GET /api/v1/admin/users`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [
    {
      "id": 1,
      "email": "user@example.com",
      "username": "johndoe",
      "email_verified": false,
      "is_active": true,
      "last_login": "2024-01-15T10:30:00Z",
      "roles": [
        {
          "id": 1,
          "role_name": "user",
          "description": "Standard user role"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 42,
    "total_pages": 5
  }
}
```

**Error Responses:**
- `401` - Unauthorized
- `403` - Forbidden (not admin)

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/users?page=1&limit=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### Assign Roles

Assigns roles to a user (replaces existing roles).

**Endpoint:** `PUT /api/v1/admin/users/:id/roles`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**URL Parameters:**
- `id`: User ID

**Request Body:**
```json
{
  "role_ids": [1, 2]
}
```

**Notes:**
- Removes all existing roles and assigns new ones
- Role IDs: 1 (user), 2 (admin), 3 (moderator)
- Must provide at least one role

**Success Response (200):**
```json
{
  "success": true,
  "message": "Roles assigned successfully",
  "data": null
}
```

**Error Responses:**
- `400` - Invalid user ID or role ID
- `401` - Unauthorized
- `403` - Forbidden (not admin)
- `404` - User or role not found

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/5/roles \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "role_ids": [2]
  }'
```

---

### Deactivate User

Marks a user account as inactive.

**Endpoint:** `DELETE /api/v1/admin/users/:id`

**Headers:**
```
Authorization: Bearer <access_token>
```

**URL Parameters:**
- `id`: User ID

**Success Response (200):**
```json
{
  "success": true,
  "message": "User deactivated successfully",
  "data": null
}
```

**Error Responses:**
- `400` - Invalid user ID
- `401` - Unauthorized
- `403` - Forbidden (not admin)
- `404` - User not found

**cURL Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/5 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Health Check

### Check API Health

Checks the health of the API and database connection.

**Endpoint:** `GET /health`

**No authentication required**

**Success Response (200):**
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": 1705318200
}
```

**Unhealthy Response (503):**
```json
{
  "status": "unhealthy",
  "database": "disconnected",
  "error": "connection refused",
  "timestamp": 1705318200
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/health
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error |
| 503 | Service Unavailable - Health check failed |

---

## Rate Limiting

(To be implemented with Redis)

---

## Pagination

List endpoints support pagination with the following query parameters:

- `page`: Page number (starts at 1)
- `limit`: Number of items per page (default: 10, max: 100)

Response includes pagination metadata:
```json
{
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 42,
    "total_pages": 5
  }
}
```

---

## Authentication Flow Diagram

```
1. Register/Login
   ↓
2. Receive access_token (24h) + refresh_token (7 days)
   ↓
3. Use access_token in Authorization header
   ↓
4. When access_token expires:
   ↓
5. Call /auth/refresh with refresh_token
   ↓
6. Receive new access_token + refresh_token
```

---

## Postman Collection

Import the following base URL and authentication settings:

**Base URL:** `http://localhost:8080/api/v1`

**Authorization (for protected routes):**
- Type: Bearer Token
- Token: `{{access_token}}`

**Environment Variables:**
- `base_url`: `http://localhost:8080`
- `access_token`: (set after login)
- `refresh_token`: (set after login)

---

For more information, see the [README.md](README.md).
