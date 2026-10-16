# RSS Aggregator

A Go-based backend service that collects and aggregates content from multiple RSS feeds. Users can register, add feeds, follow feeds from other users, and retrieve the latest posts from their subscriptions — all backed by PostgreSQL.

---

## Features

- **User registration** with auto-generated API keys
- **Feed management** — add and list RSS feeds
- **Feed follows** — follow/unfollow feeds across users
- **Automated scraping** — background goroutine fetches up to N feeds concurrently every minute
- **Post retrieval** — paginated list of posts from followed feeds, sorted newest-first
- **Graceful shutdown** — SIGINT/SIGTERM drains in-flight scrapes and closes the HTTP server cleanly
- **Request logging** — every request is logged with method, path, status code, and duration

---

## Architecture

```
┌─────────────────────────────────────────┐
│              HTTP Server (chi)           │
│  /v1/healthz  /v1/users  /v1/feeds      │
│  /v1/posts    /v1/feed_follows          │
│         ↕ loggingMiddleware             │
│         ↕ middlewareAuth (ApiKey)       │
└────────────────────┬────────────────────┘
                     │
           ┌─────────▼──────────┐
           │   apiConfig.DB     │
           │  (sqlc Queries)    │
           └─────────┬──────────┘
                     │
           ┌─────────▼──────────┐
           │     PostgreSQL      │
           │  users, feeds,      │
           │  feed_follows,      │
           │  posts              │
           └────────────────────┘
                     ▲
           ┌─────────┴──────────┐
           │  startScraper()    │
           │  (goroutine, stop  │
           │   channel)         │
           └────────────────────┘
```

---

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- [goose](https://github.com/pressly/goose) (for migrations)
- [sqlc](https://sqlc.dev) (only needed to regenerate queries)

---

## Setup

### 1. Clone & configure environment

```bash
git clone https://github.com/SuhasKanwar/RSS-Aggregator
cd RSS-Aggregator
cp .env.example .env
# Edit .env with your PORT and DATABASE_URL
```

### 2. Run database migrations

```bash
goose postgres "postgresql://user:password@localhost:5432/rssaggr" up \
  -dir sql/schema
```

### 3. Build and run

```bash
go build -o RSS-Aggregator .
./RSS-Aggregator
```

Or in one step:

```bash
go run .
```

---

## API Reference

All endpoints are prefixed with `/v1`. Authenticated routes require:

```
Authorization: ApiKey <your_api_key>
```

### Health

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/v1/healthz` | No | Returns `{}` with 200 |

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/v1/users` | No | Create a new user |
| GET | `/v1/users` | Yes | Get the authenticated user |

**POST /v1/users** body:
```json
{ "name": "Alice" }
```

### Feeds

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/v1/feeds` | Yes | Add a new RSS feed |
| GET | `/v1/feeds` | No | List all feeds |

**POST /v1/feeds** body:
```json
{ "name": "Go Blog", "url": "https://go.dev/blog/feed.atom" }
```

### Feed Follows

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/v1/feed_follows` | Yes | Follow a feed |
| GET | `/v1/feed_follows` | Yes | List your followed feeds |
| DELETE | `/v1/feed_follows/{id}` | Yes | Unfollow a feed |

**POST /v1/feed_follows** body:
```json
{ "feed_id": "<uuid>" }
```

### Posts

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/v1/posts` | Yes | Get posts from followed feeds |

**Query params:**
- `limit` — number of posts (1–100, default 10)

---

## Running Tests

```bash
go test ./... -v
```

17 tests covering: auth middleware, date parsing (8 formats), handler validation, and logging middleware.

---

## Rollback Guide

Each improvement phase is tagged for easy rollback:

| Tag | What it contains |
|-----|-----------------|
| `v0.1-bugfixes` | All critical bug fixes + input validation + pagination |
| `v0.2-hardening` | Graceful shutdown, DB pool, timeouts, 404/405 handlers |
| `v0.3-tests-and-logging` | Request logging middleware + full test suite |

To roll back to a specific tag:
```bash
git checkout v0.1-bugfixes
```

---

## Project Structure

```
.
├── main.go                    # Entry point, router, graceful shutdown
├── rss.go                     # RSS XML parsing + shared HTTP client
├── scraper.go                 # Background scraper goroutine
├── models.go                  # DB → API type conversions
├── json.go                    # respondWithJSON / respondWithError helpers
├── middleware_auth.go         # ApiKey authentication middleware
├── middleware_logging.go      # Request logging middleware
├── user_controller.go         # User handlers
├── feed_controller.go         # Feed handlers
├── feed_follows_controller.go # Feed follow handlers
├── readiness_controller.go    # Health check
├── handlers_test.go           # Handler integration tests
├── scraper_test.go            # parsePubDate unit tests
├── internal/
│   ├── auth/                  # API key extraction
│   │   └── auth_test.go
│   └── database/              # sqlc-generated query code
├── sql/
│   ├── queries/               # SQL source for sqlc
│   └── schema/                # goose migration files
└── sqlc.yaml                  # sqlc configuration
```
