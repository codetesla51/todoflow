# TodoAPI - Production Patterns Playground

[![Demo](https://img.shields.io/badge/Demo-Live-green)](https://todoflow-black.vercel.app/)

A Todo API that I built to experiment with production-grade backend architecture. The Todo domain is intentionally simple - I wanted to focus on implementing distributed systems patterns without getting lost in business logic.

**What I'm actually demonstrating here:**
- Distributed rate limiting with automatic Redis → PostgreSQL failover
- Multi-layer caching with pattern-based invalidation
- Fault-tolerant design (fail-open strategies)
- High-availability infrastructure patterns

This isn't trying to be the best Todo app. It's me learning how to build systems that don't fall over under load.

---

## Why This Exists

Most Todo apps are basic CRUD tutorials. I wanted to build something that explores real backend engineering problems:

- What happens when your rate limiter's Redis instance goes down?
- How do you invalidate cache entries across paginated results without breaking everything?
- How do you design APIs that stay available even when dependencies fail?

The Todo domain is simple enough that I could focus entirely on the infrastructure and patterns, not complex business logic.

---

## Tech Stack

- **Backend**: Go (Gin framework)
- **Frontend**: Svelte 5 + Tailwind CSS 4
- **Database**: PostgreSQL
- **Cache**: Redis
- **Deployment**: Docker Compose

---

## The Interesting Parts

### 1. Fault-Tolerant Rate Limiting

I'm using my own [limitz](https://github.com/codetesla51/limitz) library here. The rate limiting has two layers:

**IP-Based (Public Routes)**
- Token Bucket algorithm: 100 burst capacity, refills at 10/sec
- Prevents brute-force attacks on auth endpoints
- Applied to `/auth/*` and `/ping`

**User-Based (Authenticated Routes)**
- Sliding Window Counter: 1000 requests per hour
- More precise quota management
- Applied to all `/api/*` endpoints

**The Failover Strategy**

Here's where it gets interesting. The rate limiter tries Redis first (fast, distributed), but if Redis goes down, it automatically falls back to PostgreSQL. If both fail, the middleware fails-open to keep the API available rather than blocking all traffic.

Primary: Redis (in-memory, distributed)  
↓ (if unavailable)  
Secondary: PostgreSQL (slower but reliable)  
↓ (if both fail)  
Fail-Open: Allow requests through

This means the API stays up even when infrastructure has issues.

---

### 2. Cache Strategy (The Tricky Part)

To reduce database pressure, I implemented a Cache-Aside pattern with Redis for `/api/todos` and `/api/profile`.

**How it works:**

**Read Flow:**
1. Check Redis first
2. Cache hit? Return immediately
3. Cache miss? Query PostgreSQL, populate Redis, then return

**Write Flow (The Interesting Bit):**

When a user creates, updates, or deletes a Todo, I don't try to update the cache. Instead, I invalidate it using pattern-based deletion.

**Paginated Cache Keys**

Cache keys look like: `todos:user:{id}:limit:{10}:page:{1}`

This prevents a common bug where you cache page 1, then someone requests page 2, and you return stale page 1 data.

**Pattern-Based Invalidation**

When any Todo changes, I run a Redis `SCAN` for `todos:user:{id}:*` and delete all matching keys. This invalidates every cached page for that user in one operation.

Why not just update the cache? Because with pagination, you'd need to recalculate every page. It's simpler and more reliable to just nuke everything and let the next read rebuild it.

---

### 3. Backend Architecture

I structured the backend with clear separation:

```
handlers/     → HTTP request/response logic
models/       → Data structures and GORM schemas
middleware/   → Auth, rate limiting, CORS
services/     → External infrastructure (Redis, etc)
```

Nothing groundbreaking, just keeping things organized so I can find stuff later.

---

## API Endpoints

### Authentication
- `POST /auth/register` - Create account
- `POST /auth/login` - Get JWT token

### Todo Operations (Requires JWT)
- `GET /api/profile` - Get user info
- `GET /api/todos?limit=10&page=1` - List todos (paginated)
- `POST /api/todos` - Create todo
- `GET /api/todos/:id` - Get specific todo
- `PUT /api/todos/:id` - Update todo
- `PATCH /api/todos/:id/status` - Toggle completed/pending
- `DELETE /api/todos/:id` - Delete todo

---

## Running It Locally

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (optional, for development)
- Node.js 20+ (optional, for frontend development)

### Quick Start

```bash
docker compose up --build
```

This spins up PostgreSQL, Redis, the Go backend, and the Svelte frontend.

### Environment Variables

Create a `.env` file:

```env
DATABASE_URL=postgresql://user:password@localhost:5432/todos
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key-here
PUBLIC_API_URL=http://localhost:8080
```

---

## Performance Notes

I haven't done extensive load testing yet, but here's what I've observed:

- **Rate limiter overhead**: ~2ms per request with Redis, ~8ms with PostgreSQL fallback
- **Cache effectiveness**: Reduces database queries by roughly 80% on typical usage patterns
- **Concurrent connections**: Handles 1000+ concurrent users without issues (tested with `hey`)

These aren't scientific benchmarks, just rough observations from local testing.

---

## What I Learned Building This

**Cache invalidation is harder than I thought**

Getting pagination + caching right took several iterations. My first attempt cached entire result sets, which broke when users had different page sizes. Pattern-based invalidation turned out to be way more reliable than trying to be clever with TTLs.

**Fail-open vs fail-closed is a real decision**

I chose to fail-open on rate limiting (allow traffic when systems are down) because availability mattered more than strict quotas for this project. In a different context (say, protecting a payment API), I'd probably fail-closed.

**Testing distributed systems locally is awkward**

Docker Compose helped, but simulating Redis failures and failover scenarios was still pretty manual. I'd probably add chaos engineering tools if I were doing this in production.

**Graceful shutdown matters more than I expected**

Implementing proper graceful shutdown with connection draining prevented a lot of weird edge cases during restarts. Letting in-flight requests finish before shutting down the server turned out to be way more important than I initially thought.

---

## Tech Stack Details

**Backend:**
- Go 1.25
- Gin (HTTP framework)
- GORM (ORM)
- JWT for authentication
- Custom [limitz](https://github.com/codetesla51/limitz) library for rate limiting

**Frontend:**
- Svelte 5
- Tailwind CSS 4

**Infrastructure:**
- PostgreSQL (primary database)
- Redis (caching + rate limiting)
- Docker Compose (local orchestration)

---