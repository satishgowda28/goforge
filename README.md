# GoForge

Self-hosted AI agent orchestration engine written in Go. Implements the **ReAct (Reason + Act)** loop — an agent receives a task, reasons about which tool to use, executes it, observes the result, and repeats until the task is complete.

---

## What It Is

GoForge is a backend engine, not a chat app or a thin LLM wrapper. It lets you define agents in YAML, give them tools, and run them against tasks via a REST API with real-time SSE streaming.

---

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go |
| HTTP Router | `chi` |
| LLM Providers | Anthropic Claude (primary), OpenAI (secondary) |
| Database | Postgres |
| DB Access | `sqlc` |
| Config | `viper` |
| CLI | `cobra` |
| Streaming | SSE (Server-Sent Events) |
| Concurrency | Goroutines + channels + `sync` |

---

## Project Structure

```
goforge/
├── cmd/goforge/main.go              # Entry point
├── internal/
│   ├── agent/                       # ReAct loop, memory
│   ├── api/                         # chi router, handlers, SSE
│   ├── llm/                         # LLM provider interface + clients
│   ├── tools/                       # Tool interface, registry, built-in tools
│   ├── worker/                      # Bounded goroutine worker pool
│   └── db/                          # sqlc queries and generated code
├── pkg/config/config.go             # App-wide config structs
├── agents/                          # YAML agent definitions
├── migrations/                      # SQL migration files
├── config.yml                       # App configuration
└── sqlc.yaml                        # sqlc config
```

---

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/agents` | List agents from YAML configs |
| `POST` | `/agents/run` | Start an agent run |
| `GET` | `/runs` | List recent runs |
| `GET` | `/runs/:id` | Get run with all steps |
| `GET` | `/runs/:id/stream` | SSE stream of live agent events |
| `DELETE` | `/runs/:id` | Cancel an in-progress run |

**SSE event types:** `thought`, `action`, `observation`, `done`, `error`

---

## Agent Config (YAML)

Define agents in `agents/*.yaml` — no code changes needed.

```yaml
name: researcher
description: Researches topics and saves summaries
system_prompt: |
  You are a research assistant. Search the web, synthesize findings,
  and write a structured summary to a file. Cite your sources.
tools:
  - http_fetch
  - file_write
max_steps: 10
memory_type: short_term   # short_term | long_term | none
```

---

## Configuration

`config.yml`:

```yaml
server:
  port: 8080
  host: localhost

llm:
  provider: anthropic
  model: claude-sonnet-4-5
  timeout_seconds: 30

tools:
  working_dir: ./workspace
  shell_allowed_commands: [ls, cat, grep, find]

worker:
  pool_size: 5

db:
  db_url: postgres://localhost:5432/goforge
```

Override any value with environment variables:

```bash
GOFORGE_LLM_API_KEY=sk-...
```

---

## Getting Started

```bash
# Install dependencies
go mod tidy

# Copy env file
cp .example_env .env
# Add your API key to .env

# Run the server
go run ./cmd/goforge

# Health check
curl http://localhost:8080/health
```

---

## CLI

```bash
goforge run --agent researcher --task "Summarize Go 1.24 release notes"
goforge run --agent researcher --task "..." --stream
goforge runs list
goforge runs get <run-id>
goforge runs cancel <run-id>
goforge agents list
```

---

## The ReAct Loop

```
1. Receive task
2. Build prompt: system prompt + task + memory + tool descriptions
3. Call LLM → reasoning + tool choice + tool input
4. Execute tool via worker pool
5. Append observation to context
6. Repeat from step 3 until:
   - LLM returns final answer (no tool call)
   - max_steps reached
   - Context cancelled
```

---

## Development

```bash
# Run tests
go test ./...

# Run with race detector (important for concurrency)
go test -race ./...

# Generate sqlc code (after editing .sql files)
sqlc generate

# Build binary
go build -o goforge ./cmd/goforge
```

---

## Implementation Phases

| Phase | Focus |
|---|---|
| 1 — Foundation | Config, LLM client, basic REST API |
| 2 — Agent + Tools | Tool interface, ReAct loop, YAML agents |
| 3 — Concurrency + Streaming | Worker pool, goroutines, channels, SSE |
| 4 — Persistence + CLI | Postgres, sqlc, memory, Cobra CLI |
