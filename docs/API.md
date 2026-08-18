# FTN-AI API Reference

## Base URL
```
https://api.ftn-ai.com/api/v1
```

## Authentication

All API requests require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Endpoints

### Authentication

#### POST /auth/login
Login with email and password.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "admin"
  },
  "expires_at": "2024-08-19T10:30:00Z"
}
```

#### POST /auth/logout
Logout and invalidate token.

**Response:**
```json
{
  "message": "logged out successfully"
}
```

#### POST /auth/refresh
Refresh authentication token.

**Request:**
```json
{
  "token": "old-token-here"
}
```

### Users

#### GET /users
List all users (admin only).

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "admin",
    "active": true,
    "created_at": "2024-08-18T10:00:00Z"
  }
]
```

#### GET /users/{id}
Get user by ID.

#### POST /users
Create new user (admin only).

**Request:**
```json
{
  "email": "newuser@example.com",
  "name": "Jane Doe",
  "password": "securepassword",
  "role": "user"
}
```

#### PUT /users/{id}
Update user.

#### DELETE /users/{id}
Delete user.

### Services

#### GET /services
List all services.

#### GET /services/{id}
Get service details.

#### POST /services
Create new service.

#### PUT /services/{id}
Update service.

#### DELETE /services/{id}
Delete service.

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request parameters",
  "details": ["email is required"]
}
```

### 401 Unauthorized
```json
{
  "error": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": "You don't have permission to access this resource"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Rate Limiting

API requests are limited to:
- 100 requests per minute for authenticated users
- 10 requests per minute for unauthenticated requests

Rate limit info in response headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1629374400
```