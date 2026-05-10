# Daily Stand-Up System (Go)

A backend application for managing Daily Stand-Up sessions, tasks, sub-tasks, notes, and blockers. This project is a Go implementation of a system originally designed in Kotlin/Spring, utilizing modern Go idioms and a clean architecture.

## Features

- **Task Management**: Create, update, start, complete, and cancel tasks.
- **Sub-Task Management**: Break down tasks into smaller units of work.
- **Blockers**: Track issues that prevent progress on tasks or sub-tasks.
- **Notes**: Categorized notes (General, Learning, Links, Ideas, etc.) for tasks and sub-tasks.
- **Daily Sessions**: Log daily stand-up summaries and raw notes.
- **Business Logic Enforcement**: 
    - Automated status transitions (e.g., tasks are automatically blocked if they have unresolved blockers).
    - Policy guards (e.g., cannot complete a task with incomplete sub-tasks or unresolved blockers).
    - Ownership validation (e.g., ensuring sub-tasks and notes belong to the correct parent task).

## Tech Stack

- **Language**: Go 1.25+
- **Web Framework**: [Gin](https://gin-gonic.com/)
- **Database**: PostgreSQL
- **Database Driver**: [lib/pq](https://github.com/lib/pq)
- **SQL Code Generation**: [sqlc](https://sqlc.dev/)
- **Configuration**: [godotenv](https://github.com/joho/godotenv)

## Project Structure

```text
.
├── db
│   ├── queries          # SQL query definitions for sqlc
│   └── sqlc             # Generated Go code from SQL queries
├── internal
│   ├── config           # Configuration and environment variables
│   ├── handler          # Gin HTTP handlers
│   ├── models           # Domain models and DTOs
│   ├── repository       # Database repository interfaces and implementations
│   └── service          # Business logic and services
├── .env.example         # Template for environment variables
├── go.mod               # Go module definition
├── init.sql             # Database schema initialization
├── main.go              # Application entry point
└── sqlc.yaml            # sqlc configuration
```

## Getting Started

### Prerequisites

- Go 1.25 or later
- PostgreSQL
- `sqlc` (optional, for regenerating database code)

### Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd daily-app-go
   ```

2. **Configure Environment Variables**:
   Copy the example environment file and update it with your database credentials.
   ```bash
   cp .env.example .env
   ```
   Edit `.env`:
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/daily_standup?sslmode=disable
   SERVER_PORT=8080
   ```

3. **Initialize the Database**:
   Run the `init.sql` script against your PostgreSQL instance:
   ```bash
   psql -d daily_standup -f init.sql
   ```

4. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

### Running the Application

To start the server:
```bash
go run main.go
```
The server will start on the port defined in your `.env` file (default: 8080).

## API Endpoints

### Tasks
- `GET /tasks` - List all tasks
- `GET /tasks/:id` - Get task by ID
- `GET /tasks/:id/details` - Get task with all its sub-tasks, notes, and blockers
- `POST /tasks` - Create a new task
- `PATCH /tasks/:id` - Update task title
- `POST /tasks/:id/start` - Set status to IN_PROGRESS
- `POST /tasks/:id/complete` - Set status to DONE
- `POST /tasks/:id/reopen` - Set status to TODO
- `POST /tasks/:id/cancel` - Set status to CANCELED

### Sub-Tasks
- `GET /tasks/:taskId/subtasks` - List sub-tasks for a task
- `POST /tasks/:taskId/subtasks` - Create a sub-task
- `PATCH /tasks/:taskId/subtasks/:id` - Update sub-task
- `POST /tasks/:taskId/subtasks/:id/complete` - Mark as completed
- `POST /tasks/:taskId/subtasks/:id/reopen` - Mark as incomplete

### Blockers
- `GET /tasks/:taskId/blockers` - List blockers for a task
- `POST /tasks/:taskId/blockers` - Create a blocker
- `POST /tasks/:taskId/blockers/:id/resolve` - Resolve a blocker

### Notes
- `GET /tasks/:taskId/notes` - List notes for a task
- `POST /tasks/:taskId/notes` - Create a note

### Daily Sessions
- `GET /sessions` - List all daily sessions
- `POST /sessions` - Log a new daily session

## Development

### Regenerating Database Code
If you modify the SQL queries in `db/queries/*.sql`, regenerate the Go code using `sqlc`:
```bash
sqlc generate
```
