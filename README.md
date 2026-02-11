# TodoAPI Full-Stack Application

A production-ready Todo management system built with a high-performance Go backend and a Svelte frontend. This project implements advanced rate limiting with distributed state, fail-over mechanisms, and multi-layer caching.

## Architecture Overview

The system is designed as a decoupled full-stack application leveraging the following stack:

- **Backend**: Go (Gin Gonic)
- **Frontend**: Svelte 5 / Tailwind CSS 4
- **Database**: PostgreSQL (GORM)
- **Cache**: Redis
- **Infrastructure**: Docker Compose for localized orchestration

### Backend Design Patterns

The backend follows a service-oriented architecture with clear separation of concerns:
- **Handlers**: Manage HTTP request/response cycle.
- **Models**: Define data structures and GORM schemas.
- **Middleware**: Handle cross-cutting concerns (Auth, Rate Limiting, CORS).
- **Services**: Abstract external infrastructure like Redis connectivity.

## Rate Limiting Implementation

The application uses the [limitz](https://github.com/codetesla51/limitz) library for sophisticated traffic management. It implements a multi-tier, fault-tolerant rate limiting strategy.

### 1. Dual-Layer Throttling
- **IP-Based Limiting**: Applied to public routes (Auth, Ping) to prevent brute-force attacks and general DOS. Uses a **Token Bucket** algorithm (100 burst, 10/sec refill).
- **User-Based Limiting**: Applied to authenticated API routes. Uses a **Sliding Window Counter** algorithm (1000 requests per hour) for precise quota management.

### 2. High-Availability Fail-Over
The rate limiter is configured to be fault-tolerant:
- **Primary**: Redis Store (Distributed, In-memory).
- **Secondary**: PostgreSQL Store (Fail-over).
- **Logic**: If the Redis pool becomes unavailable, the system automatically transitions to using the database for rate limit tracking. If both systems fail, the middleware fails-open to ensure service availability.

## Optimized Caching Strategy

The system utilizes a **Cache-Aside (Lazy Loading)** strategy with Redis for the `/api/todos` and `/api/profile` endpoints to minimize database load.

### 1. Cache-Aside Workflow
- **Read**: The application first checks Redis for the requested data. If present (Cache Hit), it returns immediately. If absent (Cache Miss), it fetches from PostgreSQL, populates the cache for future requests, and then returns.
- **Write**: When data is modified (Create/Update/Delete), the corresponding cache entries are invalidated (deleted) rather than updated, ensuring that the next read operation fetches the most recent data from the database.

### 2. Paginated Cache Keys
- Cache keys are generated using a specific pattern: `todos:user:{id}:limit:{n}:page:{n}`. This ensures that different pagination views do not return incorrect cached data.

### 3. Pattern-Based Invalidation
- When a user creates, updates, or deletes a Todo, the system performs a `SCAN` for `todos:user:{id}:*` to invalidate all relevant cache entries across all pages simultaneously. This maintain cache consistency with minimal overhead.

## API Documentation

### Authentication
- `POST /auth/register`: Create a new account.
- `POST /auth/login`: Authenticate and receive a JWT.

### Secure API (Required: Bearer Token)
- `GET /api/profile`: Retrieve user information and todos.
- `GET /api/todos`: List todos (Supports `?limit=n&page=n`).
- `POST /api/todos`: Create a new todo.
- `GET /api/todos/:id`: Retrieve specific todo.
- `PUT /api/todos/:id`: Update todo content.
- `PATCH /api/todos/:id/status`: Toggle todo status (pending/completed).
- `DELETE /api/todos/:id`: Remove a todo.

## Local Development

### Prerequisites
- Docker and Docker Compose
- Go 1.25+ (for local development)
- Node.js 20+ (for local development)

### Deployment with Docker
To spin up the entire stack including databases:
```bash
docker compose up --build
```

### Environment Configuration
The following environment variables are required:
- `DATABASE_URL`: Connection string for PostgreSQL.
- `REDIS_HOST`: Hostname for the Redis server.
- `REDIS_PORT`: Port for the Redis server.
- `JWT_SECRET`: Secret key for signing tokens.
- `PUBLIC_API_URL`: (Frontend) The base URL of the backend API.

