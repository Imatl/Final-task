# SupportFlow AI

AI-powered customer support operations assistant.

## Project Structure

```
src/           — Go backend (module: supportflow)
web/           — React frontend (Vite + TypeScript)
config/        — .properties config files
```

## Running

```bash
# Backend (from src/)
go run .

# Frontend (from web/)
npm run dev
```

Backend runs on :8080, frontend on :5173 with proxy to backend.

## Config

Config: `config/local-application.properties` (not in git — has secrets).
Copy from `.env.example` and fill in values.
Config keys override via env vars: `db.postgres.host` → `DB_POSTGRES_HOST`.

## Go Backend

- Entry: `src/main.go`
- Config parser: `src/core/config.go` (.properties format)
- API handlers: `src/api/{chat,tickets,analytics,settings}/`
- AI service: `src/services/ai/` (multi-provider: Anthropic + OpenAI)
- DB: `src/db/postgre/` (pgx pool), `src/db/redis/` (go-redis)
- Migrations: `src/db/postgre/migrations/*.sql` (auto-applied on startup)
- GOPROXY: system is set to corporate nexus, use `GOPROXY=https://proxy.golang.org,direct` for installs

## Frontend

- Stack: React 18, TypeScript, Vite, Tailwind v4, Zustand, React Query, Recharts, Lucide
- Pages: CustomerChat, AgentDashboard, Analytics, Settings
- Dark cosmic theme (velvet/neon palette from go-llm-router)
- Path alias: `@/` → `src/`

## Dependencies

- PostgreSQL (required)
- Redis (optional — server starts without it)
- Anthropic API key or OpenAI API key

## API Endpoints

- `POST /api/chat` — send message, get AI response with tool use
- `GET /api/chat/ws` — WebSocket chat
- `GET /api/tickets` — list (filter: status, priority, agent_id, category)
- `GET /api/tickets/{id}` — detail with messages, actions, customer
- `PUT /api/tickets/{id}/status` — update status
- `PUT /api/tickets/{id}/assign` — assign agent
- `POST /api/tickets/actions/approve` — approve/reject AI action
- `GET /api/analytics/overview` — stats
- `GET /api/analytics/agents` — agent performance
- `GET /api/settings/providers` — list AI providers
- `PUT /api/settings/providers` — switch provider
- `GET /api/settings/metrics` — LLM performance metrics

## Code Style

- Go: standard formatting, no comments in new code
- TypeScript: strict mode, no comments, no emojis
- Tailwind: use custom theme tokens (cosmic-*, velvet-*, neon-*)
