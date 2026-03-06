## Description

URL shortener service with a Go backend and web UI. Create short links, redirect by short code, and track visit statistics (IP, user-agent, referer). Uses PostgreSQL, Gin, Caddy (for static assets), and goose for migrations.

## Run

```bash
# Copy env and set DATABASE_URL
make generate-config

# Apply migrations
make migrate-up

# Run locally
make run-dev
```
