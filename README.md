# GTBS API Gateway

A high-performance API Gateway for the **Global Traffic Booking Service (GTBS)**, built with [Go](https://golang.org/) and [Fiber](https://gofiber.io/).

Acts as the single entry point for all downstream microservices — handling JWT authentication, rate limiting, reverse proxying, and per-service circuit breaking.

---

## Features

- **JWT Authentication** — validates Bearer tokens on all protected routes
- **Rate Limiting** — sliding window limiter per client IP
- **Reverse Proxy** — forwards requests to the correct downstream service
- **Circuit Breaker** — per-service breaker (Sony gobreaker) to prevent cascading failures
- **Health Check** — unauthenticated `/health` endpoint for liveness probes

---

## Project Structure

```
gtbs-api-gateway/
├── main.go                  # Entry point
├── go.mod                   # Module definition and dependencies
├── .env                     # Runtime config (not committed to Git)
├── config/
│   └── config.go            # Loads all environment variables
├── middleware/
│   ├── auth.go              # JWT validation middleware
│   └── ratelimit.go         # Rate limiting middleware
├── proxy/
│   └── proxy.go             # Circuit breaker + reverse proxy
└── router/
    └── router.go            # Route definitions and middleware wiring
```

---

## Prerequisites

- [Go 1.23+](https://golang.org/dl/)
- A running instance of each downstream service (or mocked URLs for local testing)

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/gtbs-api-gateway.git
cd gtbs-api-gateway
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

Then edit `.env`:

```bash
# Server
PORT=8080

# JWT — change this to a strong secret in production
JWT_SECRET=your-super-secret-key

# Rate Limiter
RATE_LIMIT_MAX=100        # max requests per window per IP
RATE_LIMIT_WINDOW=60      # window size in seconds

# Circuit Breaker
BREAKER_MAX_REQUESTS=3    # test requests allowed in HALF-OPEN state
BREAKER_TIMEOUT=10        # seconds before retrying after OPEN
BREAKER_FAILURES=5        # consecutive failures before tripping

# Downstream Service URLs
BOOKING_SERVICE_URL=http://localhost:8081
COMPATIBILITY_SERVICE_URL=http://localhost:8082
ROUTE_SERVICE_URL=http://localhost:8083
USER_SERVICE_URL=http://localhost:8084
AUDIT_SERVICE_URL=http://localhost:8085

# CORS (optional; set for local frontend development)
# Example: CORS_ALLOW_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
CORS_ALLOW_ORIGINS=
```

### 4. Run the Gateway

```bash
go run main.go
```

You should see:

```
┌───────────────────────────────────────────────────┐
│                   Fiber v2.x.x                    │
│               http://127.0.0.1:8080               │
└───────────────────────────────────────────────────┘
```

---

## API Routes

| Method | Path | Auth Required | Forwards To |
|--------|------|---------------|-------------|
| `GET` | `/health` | ❌ | — |
| `ALL` | `/bookings/*` | ✅ | Booking Service |
| `ALL` | `/compatibility/*` | ✅ | Journey Compatibility Service |
| `ALL` | `/routes/*` | ✅ | Route Management Service |
| `ALL` | `/users/*` | ✅ | User & Notification Service |
| `ALL` | `/audit/*` | ✅ | Audit & Observability Service |

---

## Testing

### Health Check (no token needed)

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{ "status": "ok" }
```

### Authenticated Request

First generate a test JWT — you can use [jwt.io](https://jwt.io) with your `JWT_SECRET` and HS256 algorithm.

Then:

```bash
curl http://localhost:8080/bookings/123 \
  -H "Authorization: Bearer <your-token>"
```

### Missing Token (should return 401)

```bash
curl http://localhost:8080/bookings/123
```

Expected response:
```json
{ "error": "missing auth header" }
```

### Rate Limit (should return 429 after threshold)

```bash
for i in {1..110}; do
  curl -s -o /dev/null -w "%{http_code}\n" \
    http://localhost:8080/health
done
```

You should see `200` for the first 100 requests, then `429`.

---

## Circuit Breaker Behaviour

Each downstream service has its own independent circuit breaker with three states:

| State | Behaviour |
|-------|-----------|
| **CLOSED** | Normal operation, requests forwarded |
| **OPEN** | Service returning `503 Service Unavailable` immediately |
| **HALF-OPEN** | Allows `BREAKER_MAX_REQUESTS` test requests through |

The breaker trips to OPEN after `BREAKER_FAILURES` consecutive failures, waits `BREAKER_TIMEOUT` seconds, then moves to HALF-OPEN to test recovery.

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/gofiber/fiber/v2` | HTTP framework |
| `github.com/golang-jwt/jwt/v5` | JWT parsing and validation |
| `github.com/sony/gobreaker` | Circuit breaker |
| `github.com/joho/godotenv` | `.env` file loading |

---

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Port the gateway listens on |
| `JWT_SECRET` | — | Secret key for JWT verification |
| `RATE_LIMIT_MAX` | `100` | Max requests per window per IP |
| `RATE_LIMIT_WINDOW` | `60` | Rate limit window in seconds |
| `BREAKER_MAX_REQUESTS` | `3` | HALF-OPEN test request count |
| `BREAKER_TIMEOUT` | `10` | Seconds before HALF-OPEN retry |
| `BREAKER_FAILURES` | `5` | Failures before breaker trips |
| `BOOKING_SERVICE_URL` | — | Booking Service base URL |
| `COMPATIBILITY_SERVICE_URL` | — | Compatibility Service base URL |
| `ROUTE_SERVICE_URL` | — | Route Management Service base URL |
| `USER_SERVICE_URL` | — | User & Notification Service base URL |
| `AUDIT_SERVICE_URL` | — | Audit & Observability Service base URL |
| `CORS_ALLOW_ORIGINS` | — | Comma-separated allowed CORS origins; when unset, CORS middleware is disabled |

---

## Notes

- `.env` is excluded from version control via `.gitignore` — never commit secrets
- All services are independently deployable; update the `*_SERVICE_URL` variables to point to wherever each service is running
- For Kubernetes deployment, replace `.env` values with ConfigMaps and Secrets
```

***

Also create a `.env.example` file alongside it (this one **is** safe to commit — no real secrets):

```bash
PORT=8080
JWT_SECRET=
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=60
BREAKER_MAX_REQUESTS=3
BREAKER_TIMEOUT=10
BREAKER_FAILURES=5
BOOKING_SERVICE_URL=http://localhost:8081
COMPATIBILITY_SERVICE_URL=http://localhost:8082
ROUTE_SERVICE_URL=http://localhost:8083
USER_SERVICE_URL=http://localhost:8084
AUDIT_SERVICE_URL=http://localhost:8085
CORS_ALLOW_ORIGINS=
```


# Build and start everything
docker compose up --build

# Test health check
curl http://localhost:8080/health

# Test proxy with a valid JWT
curl http://localhost:8080/bookings/anything \
  -H "Authorization: Bearer <your-token>"

# Run in background
docker compose up --build -d

# Stop everything
docker compose down
