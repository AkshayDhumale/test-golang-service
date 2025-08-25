# Go Service with Redis and PostgreSQL

A complete Go web service that demonstrates integration with Redis (for caching) and PostgreSQL (for persistent storage). This service provides a RESTful API for user management with caching capabilities.

## Features

- **RESTful API** for user CRUD operations
- **PostgreSQL** for persistent data storage
- **Redis** for caching with automatic cache invalidation
- **Docker Compose** setup for easy deployment
- **Health check** endpoint for monitoring
- **Graceful error handling** and logging
- **Environment-based configuration**

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Service    │    │     Redis       │    │   PostgreSQL    │
│                 │────│   (Cache)       │    │   (Database)    │
│  - REST API     │    │                 │    │                 │
│  - Business     │    └─────────────────┘    └─────────────────┘
│    Logic        │
└─────────────────┘
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check for the service |
| POST | `/api/v1/users` | Create a new user |
| GET | `/api/v1/users` | Get all users |
| GET | `/api/v1/users/{id}` | Get user by ID (cached) |
| PUT | `/api/v1/users/{id}` | Update user by ID |
| DELETE | `/api/v1/users/{id}` | Delete user by ID |

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone and navigate to the project directory**

2. **Start all services:**
   ```bash
   make docker-up
   ```

3. **Test the API:**
   ```bash
   make api-test
   ```

4. **View logs:**
   ```bash
   make docker-logs
   ```

5. **Stop services:**
   ```bash
   make docker-down
   ```

### Manual Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Start PostgreSQL and Redis locally**

3. **Set environment variables:**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=userdb
   export REDIS_HOST=localhost
   export REDIS_PORT=6379
   ```

4. **Run the service:**
   ```bash
   make run
   ```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | PostgreSQL username |
| `DB_PASSWORD` | postgres | PostgreSQL password |
| `DB_NAME` | userdb | PostgreSQL database name |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `REDIS_PASSWORD` | "" | Redis password |
| `PORT` | 8080 | Service port |

## Docker Services

The Docker Compose setup includes:

- **app**: Go service (port 8080)
- **postgres**: PostgreSQL database (port 5432)
- **redis**: Redis cache (port 6379)
- **pgadmin**: PostgreSQL web UI (port 8082) - Optional
- **redis-commander**: Redis web UI (port 8081) - Optional

### Access Web UIs

- **API Service**: http://localhost:8080
- **pgAdmin**: http://localhost:8082 (admin@example.com / admin)
- **Redis Commander**: http://localhost:8081

## API Usage Examples

### Create a user:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

### Get all users:
```bash
curl http://localhost:8080/api/v1/users
```

### Get a specific user (cached):
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update a user:
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

### Delete a user:
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### Health check:
```bash
curl http://localhost:8080/health
```

## Caching Strategy

- **Cache Hit**: User data is served from Redis (faster response)
- **Cache Miss**: Data is fetched from PostgreSQL and cached in Redis
- **Cache TTL**: 1 hour
- **Cache Invalidation**: Automatic on user updates and deletions
- **Cache Headers**: `X-Cache: HIT/MISS` indicates cache status

## Development

### Project Structure
```
.
├── main.go           # Main application code
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── Dockerfile       # Docker image definition
├── docker-compose.yml # Multi-service setup
├── init.sql         # Database initialization
├── Makefile         # Development commands
└── README.md        # This file
```

### Available Make Commands
```bash
make help        # Show available commands
make build       # Build the application
make run         # Run locally
make test        # Run tests
make docker-up   # Start with Docker Compose
make docker-down # Stop services
make docker-logs # View logs
make api-test    # Test API endpoints
```

## Production Considerations

1. **Environment Variables**: Use secrets management for sensitive data
2. **Database Migrations**: Implement proper migration system
3. **Monitoring**: Add metrics and logging
4. **Security**: Add authentication/authorization
5. **Load Balancing**: Use multiple service instances
6. **Database Pooling**: Configure connection pooling
7. **Redis Clustering**: Use Redis Cluster for high availability
