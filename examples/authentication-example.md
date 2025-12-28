# Authentication Flow Example

This example demonstrates the complete authentication workflow including registration, login, token refresh, and logout.

## 1. User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePassword123!",
    "name": "John Doe"
  }'
```

**Response (201 Created):**
```json
{
  "id": "user-123",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

## 2. User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePassword123!"
  }'
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJlbWFpbCI6ImpvaG4uZG9lQGV4YW1wbGUuY29tIiwiZXhwIjoxNzA1MzI2NjAwfQ.xyz",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJ0b2tlbl9pZCI6InRva2VuLWFiYyIsImV4cCI6MTcwNTkzMTQwMH0.abc",
  "expiresIn": 900,
  "tokenType": "Bearer"
}
```

**Token Details:**
- `access_token`: Valid for 15 minutes, used for API authentication
- `refresh_token`: Valid for 7 days, used to obtain new access tokens
- `expiresIn`: Access token expiration time in seconds

## 3. Making Authenticated Requests

Store the access token and use it in the Authorization header:

```bash
export ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get user profile
curl http://localhost:8080/api/v1/users/user-123 \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response (200 OK):**
```json
{
  "id": "user-123",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "tenantId": "tenant-456",
  "roles": ["user"],
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

## 4. Handling Token Expiration

When the access token expires (after 15 minutes), API requests will return 401 Unauthorized:

```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token",
  "correlationId": "abc-123-def",
  "timestamp": "2024-01-15T10:45:00Z"
}
```

## 5. Refreshing Access Token

Use the refresh token to obtain a new access token:

```bash
export REFRESH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new_token.xyz",
  "expiresIn": 900,
  "tokenType": "Bearer"
}
```

Update your stored access token:
```bash
export ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new_token.xyz"
```

## 6. User Logout

Invalidate both access and refresh tokens:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

After logout:
- The refresh token is revoked and cannot be used
- The access token is blacklisted (until its natural expiration)
- Any subsequent requests with these tokens will return 401 Unauthorized

## 7. Error Handling

### Invalid Credentials

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "WrongPassword"
  }'
```

**Response (401 Unauthorized):**
```json
{
  "error": "unauthorized",
  "message": "Invalid credentials",
  "correlationId": "abc-123-def",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Missing Authorization Header

```bash
curl http://localhost:8080/api/v1/users
```

**Response (401 Unauthorized):**
```json
{
  "error": "unauthorized",
  "message": "Missing authorization header",
  "correlationId": "abc-123-def",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Invalid Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "invalid_token"}'
```

**Response (401 Unauthorized):**
```json
{
  "error": "unauthorized",
  "message": "Invalid refresh token",
  "correlationId": "abc-123-def",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Best Practices

1. **Secure Token Storage**
   - Never store tokens in localStorage (XSS vulnerability)
   - Use httpOnly cookies or secure storage mechanisms
   - Don't log or expose tokens in error messages

2. **Token Refresh Strategy**
   - Implement automatic token refresh before expiration
   - Use interceptors in HTTP clients to handle 401 errors
   - Refresh proactively (e.g., 1 minute before expiration)

3. **Logout Handling**
   - Always call logout endpoint when user logs out
   - Clear all stored tokens
   - Redirect to login page

4. **Error Handling**
   - Implement retry logic for network errors
   - Handle 401 errors by redirecting to login
   - Display user-friendly error messages

## Example: Automatic Token Refresh (JavaScript)

```javascript
let accessToken = '';
let refreshToken = '';
let tokenExpiration = 0;

// Login
async function login(email, password) {
  const response = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  
  const data = await response.json();
  accessToken = data.access_token;
  refreshToken = data.refresh_token;
  tokenExpiration = Date.now() + (data.expiresIn * 1000);
}

// Refresh token if needed
async function ensureValidToken() {
  // Refresh 1 minute before expiration
  if (Date.now() >= (tokenExpiration - 60000)) {
    const response = await fetch('http://localhost:8080/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });
    
    const data = await response.json();
    accessToken = data.access_token;
    tokenExpiration = Date.now() + (data.expiresIn * 1000);
  }
}

// Make authenticated request
async function fetchUsers() {
  await ensureValidToken();
  
  const response = await fetch('http://localhost:8080/api/v1/users', {
    headers: { 'Authorization': `Bearer ${accessToken}` }
  });
  
  return response.json();
}

// Logout
async function logout() {
  await fetch('http://localhost:8080/api/v1/auth/logout', {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${accessToken}` }
  });
  
  accessToken = '';
  refreshToken = '';
  tokenExpiration = 0;
}
```

## Security Considerations

1. **HTTPS Only**: Always use HTTPS in production to prevent token interception
2. **Strong Passwords**: Enforce password complexity requirements
3. **Rate Limiting**: The gateway implements rate limiting to prevent brute force attacks
4. **Token Rotation**: Refresh tokens should be rotated on each use (check backend implementation)
5. **Session Management**: Implement proper session management and timeout policies
