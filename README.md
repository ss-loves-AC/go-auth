# go-auth

This project sets up a Go-based authentication app using Docker Compose with **MySQL** for data storage and **Redis** for caching.

## Services

- **MySQL**: Stores user authentication data.
- **Redis**: Caches data to improve performance.
- **Go Auth App**: Handles authentication logic.

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

1. Clone the repository:

   ```
   git clone https://github.com/ss-loves-AC/go-tube
   cd go-tube
   ```

2. Build and start the services:

   ```
   docker-compose up --build
   ```

3. Access the app at [http://localhost:8080](http://localhost:8080).

## Example curl Commands

### 1. Sign Up (POST Request)

**Request:**

```
curl -X POST http://localhost:8080/signup \
     -H "Content-Type: application/json" \
     -d '{"email": "testuser@example.com", "password": "testpassword"}'
```

**Expected Response:**

```
{
  "message": "User created successfully"
}
```

---

### 2. Sign In (POST Request)

**Request:**

```
curl -X POST http://localhost:8080/signin \
     -H "Content-Type: application/json" \
     -d '{"email": "testuser@example.com", "password": "testpassword"}'
```

**Expected Response:**

```
{
  "refresh_token": "26997b16-ed93-49d9-9ea6-4f051d9d9b88",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg2MzAwNTgsImp0aSI6IjI2OTk3YjE2LWVkOTMtNDlkOS05ZWE2LTRmMDUxZDlkOWI4OCIsInVzZXJfaWQiOjJ9.SaXDNdBJoRnG2KvD1qyteHYFQGUHvND9fseZ51dH-Qo"
}
```

---

### 3. Authorize (POST Request)

**Request:**

```
curl -X POST http://localhost:8080/authorize \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg2MzAwNTgsImp0aSI6IjI2OTk3YjE2LWVkOTMtNDlkOS05ZWE2LTRmMDUxZDlkOWI4OCIsInVzZXJfaWQiOjJ9.SaXDNdBJoRnG2KvD1qyteHYFQGUHvND9fseZ51dH-Qo"
```

**Expected Response:**

```
{
  "message": "Token is valid"
}
```

---

### 4. Revoke Token (POST Request)

**Request:**

```
curl -X POST http://localhost:8080/revoke \
     -H "Content-Type: application/json" \
     -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg2MzAwNTgsImp0aSI6IjI2OTk3YjE2LWVkOTMtNDlkOS05ZWE2LTRmMDUxZDlkOWI4OCIsInVzZXJfaWQiOjJ9.SaXDNdBJoRnG2KvD1qyteHYFQGUHvND9fseZ51dH-Qo"}'
```

**Expected Response:**

```
{
  "message": "Token revoked successfully"
}
```

---

### 5. Refresh Token (POST Request)

**Request:**

```
curl -X POST http://localhost:8080/refresh \
     -H "Content-Type: application/json" \
     -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg2MzAwNTgsImp0aSI6IjI2OTk3YjE2LWVkOTMtNDlkOS05ZWE2LTRmMDUxZDlkOWI4OCIsInVzZXJfaWQiOjJ9.SaXDNdBJoRnG2KvD1qyteHYFQGUHvND9fseZ51dH-Qo"}'
```

**Expected Response:**

```
{
  "new_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg2MzAxMDAsImp0aSI6Ijk3OTZmZWIxLWY0ZmUtNDEyOS05MTIxLTNmMDRiZWJjM2FiYiIsInVzZXJfaWQiOjJ9._30gR179TMX3GorRjhMpsI9wlY1Oz0gGMOux3pjCqUY"
}
```

---

