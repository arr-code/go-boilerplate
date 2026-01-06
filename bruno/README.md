# Bruno API Collection

This directory contains Bruno API collection files for testing the Xnoia Go Boilerplate API.

## What is Bruno?

[Bruno](https://www.usebruno.com/) is a fast, Git-friendly, open-source API client. Unlike Postman, Bruno stores collections as plain text files that work well with version control.

## Installation

Download Bruno from: https://www.usebruno.com/downloads

## Getting Started

1. **Open Bruno** and click "Open Collection"
2. **Select this folder** (`bruno/`)
3. **Select Environment**: Choose "Local" from the environment dropdown
4. **Start your API server**: `make server`
5. **Run requests** in order:
   - Health Check (verify server is running)
   - Register (creates a user and saves tokens)
   - Login (alternative to register)
   - Get Me (uses saved access_token)
   - Other endpoints

## Collection Structure

```
bruno/
├── Auth/                    # Authentication endpoints
│   ├── Register.bru        # Create new user
│   ├── Login.bru           # Login with credentials
│   ├── Refresh Token.bru   # Refresh access token
│   └── Get Me.bru          # Get current user info
├── Profile/                 # User profile endpoints
│   ├── Get Profile.bru     # Get full profile
│   ├── Update Profile.bru  # Update profile details
│   └── Upload Avatar.bru   # Upload avatar image
├── Admin/                   # Admin-only endpoints
│   ├── List Users.bru      # List all users (paginated)
│   ├── Assign Roles.bru    # Assign roles to user
│   └── Deactivate User.bru # Deactivate user account
├── environments/            # Environment configurations
│   ├── Local.bru           # Local development (localhost:8080)
│   └── Production.bru      # Production environment
├── Health Check.bru         # API health check
└── bruno.json              # Collection metadata
```

## Environment Variables

The collection uses the following variables:

- `base_url`: API base URL
- `access_token`: JWT access token (auto-set after login/register)
- `refresh_token`: JWT refresh token (auto-set after login/register)

### Local Environment
```
base_url: http://localhost:8080
```

### Production Environment
```
base_url: https://api.yourproduction.com
```

## Automatic Token Management

The Register and Login requests automatically save tokens to environment variables:

```javascript
if (res.status === 200 || res.status === 201) {
  bru.setVar("access_token", res.body.data.access_token);
  bru.setVar("refresh_token", res.body.data.refresh_token);
}
```

Protected endpoints automatically use the saved `access_token`.

## Testing Flow

### 1. Health Check
```
GET /health
```
Verify the API is running and database is connected.

### 2. Register a New User
```
POST /api/v1/auth/register
```
Creates a user with default "user" role. Tokens are saved automatically.

### 3. Get Current User
```
GET /api/v1/auth/me
```
Uses the saved access_token to get user info.

### 4. Update Profile
```
PUT /api/v1/profile
```
Update user profile details.

### 5. Admin Endpoints (Requires Admin Role)

First, assign admin role to your user:
```sql
-- Run in PostgreSQL
UPDATE user_roles SET role_id = 2 WHERE user_id = 1;
```

Or use the API after logging in with an admin account:
```
PUT /api/v1/admin/users/1/roles
{
  "role_ids": [2]
}
```

Then test admin endpoints:
- List Users
- Assign Roles
- Deactivate User

## Request Examples

### Register
```json
{
  "email": "john.doe@example.com",
  "username": "johndoe",
  "password": "SecurePass123",
  "full_name": "John Doe"
}
```

### Login
```json
{
  "login": "johndoe",
  "password": "SecurePass123"
}
```

### Update Profile
```json
{
  "full_name": "John Doe Updated",
  "phone": "+1234567890",
  "address": "123 Main Street, New York, NY",
  "date_of_birth": "1990-01-15",
  "bio": "Software engineer"
}
```

### Assign Roles
```json
{
  "role_ids": [2]
}
```

## Tests

Each request includes automated tests. After running a request, check the "Tests" tab to see:
- ✅ Status code validation
- ✅ Response structure validation
- ✅ Data type validation

Example tests:
```javascript
test("should return 200", function() {
  expect(res.status).to.equal(200);
});

test("should return access token", function() {
  expect(res.body.data.access_token).to.be.a('string');
});
```

## Tips

1. **Run in sequence**: Start with Health Check → Register → Get Me
2. **Check Tests tab**: View test results after each request
3. **Use Variables**: Access tokens are automatically saved and reused
4. **Edit User IDs**: Update user IDs in admin endpoints as needed
5. **Upload Avatar**: Change the file path in "Upload Avatar" request

## Troubleshooting

### 401 Unauthorized
- Token expired or invalid
- Run Login or Register again to get a new token

### 403 Forbidden
- Insufficient permissions
- Admin endpoints require admin role
- Assign admin role to your user first

### Connection Refused
- API server is not running
- Run: `make server`

### 404 Not Found
- Check the endpoint URL
- Ensure you're using the correct environment

## Alternative: cURL

If you prefer cURL, see [API.md](../API.md) for cURL examples.

## Learn More

- [Bruno Documentation](https://docs.usebruno.com/)
- [API Documentation](../API.md)
- [Project README](../README.md)
