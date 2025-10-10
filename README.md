# Chirpy

**Chirpy** is a small microblogging app written in Go. It provides a lightweight HTTP API for creating, reading, and deleting “chirps” (short messages), managing users, and handling JWT-based authentication. It also includes a static frontend. The project is meant as a compact, educational example.

## Table of Contents

- Features
- Prerequisites
- Build & Install
- Configuration
- Database & Migrations
- Running the Server
- CLI Commands
- API Endpoints
- Project Layout
- Development Notes
- Troubleshooting

---

## Features

- Simple HTTP API for users and chirps
- JWT authentication with refresh/revoke flows
- Polka webhook endpoint (stub)
- Static frontend at `/app/`

---

## Prerequisites

- [Go 1.20+](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- Optional tools: `psql`, `jq`, `godotenv`

---

## Build & Install

```bash
git clone https://github.com/<your-username>/chirpy.git
cd chirpy
go install
```

Or run directly:

```bash
go run . <command>
```

---

## Configuration

Chirpy loads environment variables (or a `.env` file via `godotenv`).

**Required vars:**

```
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=local
JWT_SECRET=your-jwt-secret
POLKA_KEY=your-polka-key
```

---

## Database & Migrations

Migration files live in `sql/schema/` with `-- +goose Up/Down` sections.

Apply them manually:

```bash
for f in sql/schema/*.sql; do
  awk '/-- +goose Up/{flag=1;next}/-- +goose Down/{flag=0}flag' "$f" | psql "$DB_URL"
done
```

Or with [goose](https://github.com/pressly/goose):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/schema postgres "$DB_URL" up
```

---

## Running the Server

```bash
go run .
# or
chirpy serve
```

Server runs on port **8080** by default.

---

## CLI Commands

Examples (implemented in `main.go`):

```bash
go run . addchirp "hello world"
go run . adduser "alice"
go run . serve
```

---

## API Endpoints

| Method | Endpoint            | Description          |
| ------ | ------------------- | -------------------- |
| GET    | /api/healthz        | Health check         |
| GET    | /api/chirps         | List chirps          |
| GET    | /api/chirps/{id}    | Get a chirp          |
| POST   | /api/chirps         | Create chirp         |
| DELETE | /api/chirps/{id}    | Delete chirp         |
| POST   | /api/users          | Create user          |
| PUT    | /api/users          | Update user          |
| POST   | /api/login          | Login (JWT)          |
| POST   | /api/refresh        | Refresh token        |
| POST   | /api/revoke         | Revoke refresh token |
| POST   | /api/polka/webhooks | Polka webhook        |

---

## Project Layout

- `main.go` — entrypoint & router
- `handler_*.go` — HTTP handlers
- `internal/database/` — `sqlc`-generated code
- `sql/schema/` — SQL migrations
- `assets/` — static frontend

---

## Development Notes

- Uses `sqlc` for type-safe DB access.
- If you see module errors, create a `go.work` including all required modules.

---

## Troubleshooting

| Issue                                 | Cause / Fix                                                   |
| ------------------------------------- | ------------------------------------------------------------- |
| `undefined: handlerCreateUser`        | Ensure you reference `apiCfg.handlerCreateUser` in `main.go`. |
| `pq: relation "users" does not exist` | Run DB migrations.                                            |
| Module not in workspace               | Add module path to `go.work`.                                 |
