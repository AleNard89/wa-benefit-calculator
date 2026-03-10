# RPA Benefit Calculator

A multi-tenant web application for evaluating and tracking the business benefits of Robotic Process Automation (RPA), with UiPath Orchestrator integration and an AI-powered chatbot.

Built to help RPA teams assess process automation candidates, calculate expected ROI, and monitor bot performance from a single platform.

## Features

### Process Management
- **Benefit calculation** — Multi-tab form with ~30 parameters, automatic ROI, break-even, annual savings, hours saved
- **KPI dashboard** — Each KPI has an info button explaining the formula behind the calculation
- **Process list** — Filterable by status, searchable, with linked bots column
- **Document upload** — Attach PPTX, DOCX or PDF assessment documents per process (indexed for RAG)
- **Technology multi-select** — Pick multiple technologies with a free-text "Other" option
- **Bot linking** — Associate Orchestrator bots to each process/assessment, with optional notes describing each bot's role

### UiPath Orchestrator Integration
- **Connectors** — Configure one or more UiPath Orchestrator connections with multi-folder support
- **Sync** — Pull job executions, schedules, queue definitions and queue items via OData API
- **Job monitoring** — View execution history with status, duration, errors
- **Schedules** — See active triggers with cron expressions and next occurrence
- **Queue items** — Browse queue data with status filters
- **Process-queue mapping** — Auto-detect which bots feed which queues (based on shared folders), with manual override

### AI Chat (RAG)
- **Chatbot** powered by Azure OpenAI (GPT-4.1) with Retrieval-Augmented Generation
- Context includes: uploaded documents (pgvector similarity search), process data (ROI, linked bots, notes), orchestrator data (recent executions, errors, schedules, mappings)
- Understands the relationship between assessments, bots and queues

### Platform
- **Multi-tenant** — Row-Level Security on PostgreSQL, company switcher in sidebar
- **RBAC** — Admin, Contributor, Reader roles with granular permissions
- **User management** — Full admin panel for companies, users, roles, areas
- **Soft delete** — Logical deletion with admin restore
- **Charts** — Recharts (pie, radar, bar) for visual process analysis

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.24 / Gin |
| Frontend | React 19 / TypeScript / Chakra UI v3 / Vite |
| Database | PostgreSQL 17 + pgvector |
| Cache | Redis 7 |
| AI | Azure OpenAI (GPT-4.1 + text-embedding-3-small) |
| RPA | UiPath Orchestrator (OData API) |
| Infrastructure | Docker Compose |

## Quick Start

### Prerequisites

- **Docker Desktop** installed and running
- An **Azure OpenAI** resource (optional, for chat functionality)

### Setup

```bash
# 1. Clone the repo
git clone https://github.com/AleNard89/wa-benefit-calculator.git
cd wa-benefit-calculator

# 2. Create your environment file
cp .env.example .env.dev.local

# 3. (Optional) Add your Azure OpenAI credentials in .env.dev.local

# 4. Start the stack
./dev.sh up
```

Open http://localhost:5174 and log in:
- **Email:** `admin@example.com`
- **Password:** `Admin123!`

### Commands

```bash
./dev.sh up              # Start all services
./dev.sh down            # Stop and remove containers
./dev.sh down --volumes  # Stop and wipe all data (DB, Redis)
./dev.sh restart         # Restart everything
./dev.sh logs [service]  # View logs (api, ui, postgres, redis, migrations)
./dev.sh db              # Open a psql session
./dev.sh ps              # Container status
```

### Services

| Service | URL |
|---------|-----|
| Frontend | http://localhost:5174 |
| API | http://localhost:8082/api/health |
| PostgreSQL | localhost:5435 |
| Redis | localhost:6381 |

## Project Structure

```
BenefitCalculator/
├── dev.sh                          # CLI to manage the Docker stack
├── docker-compose.base.yml         # Base services (postgres, redis)
├── docker-compose.dev.yml          # Dev overrides (hot-reload)
├── .env.example                    # Environment template
└── vertical/
    ├── api/src/                    # Go backend
    │   ├── main.go                 # Entry point, router, middleware
    │   ├── auth/                   # JWT, RBAC, brute-force protection
    │   ├── orgs/                   # Multi-tenant companies & areas
    │   ├── processes/              # CRUD, benefit calculations, document handlers
    │   ├── chat/                   # RAG, Azure OpenAI, PPTX/DOCX/PDF processor
    │   ├── orchestrator/           # UiPath connectors, sync, mapping
    │   ├── core/                   # Middleware, encryption, HTTP utils
    │   └── db/                     # PostgreSQL client
    ├── ui/src/                     # React frontend
    │   ├── Auth/                   # Login, token refresh, protected routes
    │   ├── Processes/              # Forms, list, detail, charts, KPI info
    │   ├── Orchestrator/           # Jobs, schedules, queues, mapping tab
    │   ├── Chat/                   # AI chat interface
    │   ├── Core/                   # Router, Redux, API services, Settings
    │   ├── Orgs/                   # Company management, switcher
    │   ├── Common/                 # Layout, Sidebar, ErrorBoundary
    │   └── Theme/                  # Brand colors, dark/light mode
    ├── migrations/                 # PostgreSQL migrations (sequential)
    └── media/                      # Per-company document storage
```

## Configuration

### Azure OpenAI (Chat)

Edit `.env.dev.local` with your Azure OpenAI resource:

```env
AZURE_OPENAI_ENDPOINT="https://your-resource.cognitiveservices.azure.com"
AZURE_OPENAI_API_KEY="your-api-key"
AZURE_OPENAI_CHAT_DEPLOYMENT="gpt-4.1"
AZURE_OPENAI_EMBEDDING_DEPLOYMENT="text-embedding-3-small"
AZURE_OPENAI_API_VERSION="2025-01-01-preview"
```

The chat works without these credentials but will be disabled.

### UiPath Orchestrator

Configure connectors from **Settings > Connettori** in the UI:
1. Create a connector with your Orchestrator organization name, tenant, and personal access token
2. Add the folder IDs you want to sync
3. Hit "Sync" to pull data

## License

MIT
