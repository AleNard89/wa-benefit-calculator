# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Orbita** — Web application for calculating financial benefits of RPA (Robotic Process Automation). Multi-tenant SaaS with RBAC, RAG chat (Azure OpenAI), UiPath Orchestrator integration, and Docker-based infrastructure. The UI and all user-facing text are in **Italian**.

## Commands

```bash
# Start the full stack (Docker required)
./dev.sh up

# Other useful commands
./dev.sh down            # Stop and remove containers
./dev.sh logs [service]  # View logs (api, ui, postgres, redis, migrations)
./dev.sh restart         # Restart everything
./dev.sh db              # Open psql session
./dev.sh ps              # Container status
```

No local Python, Node.js, or Go installation needed -- everything runs in Docker.

## Architecture

### Backend — Go/Gin (`vertical/api/src/`)
- **Entry point:** `main.go` — Gin router, middleware setup, DB connection
- **Auth:** JWT access/refresh tokens, RBAC with roles (Admin/Contributor/Reader), brute-force protection via Redis
- **Orgs:** Multi-tenant with company hierarchy, areas, Row-Level Security (RLS) on PostgreSQL
- **Processes:** CRUD with JSONB storage, benefit calculations (ROI, break-even, savings), linked bots, bot notes, multi-select technology
- **Chat:** RAG with Azure OpenAI (GPT-4.1 + text-embedding-3-small), pgvector, SSE streaming, PPTX/DOCX/PDF document processor. Context includes processes, orchestrator data (jobs, schedules, queues, mappings), and document chunks.
- **Orchestrator:** UiPath OData API integration — connectors, multi-folder sync, job executions, schedules, queue definitions, queue items, process-queue mapping (auto-detect + manual)

### Frontend — React/TypeScript (`vertical/ui/src/`)
- **Framework:** Vite + React + Chakra UI v3
- **State:** Redux Toolkit + RTK Query for API calls
- **Auth:** JWT with automatic token refresh, protected routes
- **Pages:** Login, Dashboard, Process List, Benefit Calculator (multi-tab form), Process Detail (with KPI info tooltips, linked bots, document upload), Orchestrator (Jobs, Schedules, Queues), Settings (Companies, Users, Roles, Connettori, Mapping Bot-Code), Chat
- **Charts:** Recharts (pie, radar, bar)
- **Guards:** `PrivateRoute` (auth), `AdminRoute` (admin/superuser only — used for Connettori and Mapping)
- **ErrorBoundary:** wraps process routes to catch and display runtime errors

### Infrastructure
- **PostgreSQL 17** with pgvector extension — main database with RLS
- **Redis 7** — session cache, brute-force protection
- **Docker Compose** — dev environment with hot-reload (reflex for Go, Vite HMR for React)

### URLs (dev)

| Service | URL |
|---------|-----|
| Frontend | http://localhost:5174 |
| API | http://localhost:8082/api/health |
| PostgreSQL | localhost:5435 |
| Redis | localhost:6381 |

## Git Workflow

**Commit after every meaningful change.** Each logical unit of work (new feature, bug fix, refactoring, security patch) must be committed immediately after verifying it compiles/passes checks. Do not accumulate uncommitted changes across multiple tasks.

- Use **Conventional Commits**: `feat:`, `fix:`, `refactor:`, `chore:`, `docs:`
- Keep commits atomic: one logical change per commit
- Always run `git diff --cached` before committing to check for secrets/credentials
- Never push without explicit user instruction

## Key Conventions

- **Italian language** throughout UI labels and user-facing text
- Brand colors: yellow accent (#FFE600), blue (#3F7DE8), green (#3ABB87), dark backgrounds
- Multi-tenant: every API request includes `X-Company-Id` header, RLS enforced at DB level
- JSONB storage for process data (flexible schema): technology (array), linkedBots (array of Orchestrator bot names), botNotes (free text for bot role descriptions)
- File storage: `vertical/media/{company-slug}/` for per-company documents (PPTX, DOCX, PDF — replace mode: new file replaces old + RAG re-index)
- Orchestrator sync: connectors with multi-folder support, auto-detect process-queue mappings during sync
- No MCP: RAG context is built via direct SQL queries (processes, orchestrator data, document chunks) and passed to Azure OpenAI as system messages
