# Short Track - High-Performance URL Shortener

Short Track is a fast, robust, and clean URL Shortener service written in Go (Golang). It leverages PostgreSQL for persistent storage and Redis for high-speed caching.

## Features

- **Base62 Encoding**: Automatically generates short, clean alphanumeric URLs from database primary keys.
- **Caching Layer**: Resolves URLs rapidly using Redis cache with a 24-hour TTL fallback to PostgreSQL.
- **Docker Compose Setup**: Spins up the application, PostgreSQL, and Redis in isolated containers with matching credentials.
- **Robust Startup**: The Go API includes built-in retry-connection loops (up to 15 retries with 2-second delays) to wait for database and cache systems to initialize properly.
- **Unit Tested**: Includes validation tests for base62 encode and decode routines.

---

## Architecture Flow

```
                      +-------------------+
                      |   Client Browser  |
                      +---------+---------+
                                |
                                | HTTP Redirect / Shorten URL
                                v
                      +-------------------+
                      |     Go API App    |
                      +----+---------+----+
                           |         |
               Cache lookup|         | Fallback & Persistence
             (Redis Get/Set)         | (PostgreSQL Select/Insert)
                           v         v
                      +----+---+ +---+----+
                      | Redis  | | Postgres|
                      +--------+ +---------+
```

---

## Getting Started

### 1. Run with Docker Compose (Recommended)

Make sure you have Docker and Docker Compose installed and running on your system.

To build and start all services (API, Postgres, Redis):
```bash
docker-compose up --build -d
```

To view logs for the API service:
```bash
docker-compose logs -f api
```

To stop all services:
```bash
docker-compose down
```

### 2. Run Locally (Without Docker)

To run the application locally, you will need:
- **Go 1.26.1+** installed
- **PostgreSQL** running locally on port `5432` with a database named `shortener`
- **Redis** running locally on port `6379`

1. Run the database migration script `/init.sql` against your PostgreSQL server to create the schema.
2. Configure your environment variables (optional, defaults are shown below):
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=pass
   export DB_NAME=shortener
   export REDIS_HOST=localhost:6379
   ```
3. Run the Go application:
   ```bash
   go run cmd/api/main.go
   ```

---

## API Reference

### 1. Shorten a URL

Generates a shortened code for a target URL.

- **URL**: `/shorten`
- **Method**: `POST`
- **Headers**: `Content-Type: application/json`
- **Request Body**:
  ```json
  {
    "url": "https://github.com"
  }
  ```
- **Response** (`200 OK`):
  ```json
  {
    "short_url": "http://localhost:8080/b"
  }
  ```

#### Example CURL request:
```bash
curl -i -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'
```

---

### 2. Redirect to Original URL

Redirects client browser to the original URL corresponding to the short code.

- **URL**: `/:code`
- **Method**: `GET`
- **Response**: `301 Moved Permanently` (Redirects to target URL)

#### Example CURL request:
```bash
curl -i http://localhost:8080/b
```

---

## Running Tests

To run the unit tests for the utility packages:
```bash
go test ./...
```
