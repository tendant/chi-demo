# chi-demo

A flexible, type-safe Go web application framework built on [chi router](https://github.com/go-chi/chi) with a powerful middleware stack system.

## Features

- ✅ **Flexible Middleware Stack** - Full control over middleware ordering, insertion, and removal
- ✅ **Functional Options Pattern** - Type-safe configuration
- ✅ **Dual Metrics Mode** - Combined (same port) or Separate (dedicated port)
- ✅ **Graceful Shutdown** - Proper signal handling
- ✅ **Structured Logging** - slog integration with environment-based configuration
- ✅ **CORS Support** - Configurable cross-origin resource sharing
- ✅ **Security Headers** - HSTS support
- ✅ **12-Factor App** - Configuration via environment variables

## Quick Start

### Simplest Example

```go
package main

import "github.com/tendant/chi-demo/app"

func main() {
    server := app.DefaultApp()
    server.Run()
}
```

### Custom Configuration

```go
myApp := app.NewApp(
    app.WithPort(8080),
    app.WithDefaultCORS(),
    app.WithMetricsCombined(),  // Metrics on same port
)

myApp.R.Get("/hello", handleHello)
myApp.Run()
```

### Metrics Modes

**Combined Mode** (simple, one port):
```go
app := app.NewApp(
    app.WithPort(3000),
    app.WithMetricsCombined(),  // Metrics at http://localhost:3000/metrics
)
```

**Separate Mode** (production, security):
```go
app := app.NewApp(
    app.WithPort(3000),
    app.WithMetricsSeparatePort(9090),  // Metrics at http://localhost:9090/metrics
)
```

### Custom Middleware Stack

```go
stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", myAuthMiddleware).
    Remove("hsts").
    Build()

myApp := app.NewApp(
    app.WithMiddlewareStack(stack),
    app.WithDefaultCORS(),
)
```

## Key Concepts

### Middleware Stack System

Control middleware ordering with a fluent builder API:

```go
stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", authMiddleware).
    Remove("hsts").
    Replace("http-logger", customLogger).
    Build()

app := app.NewApp(
    app.WithMiddlewareStack(stack),
)
```

Available methods: `Add`, `Prepend`, `Append`, `InsertBefore`, `InsertAfter`, `Replace`, `Remove`, `Enable`, `Disable`

### Configuration Options

**Application:**
- `WithConfig(AppConfig)` - Set complete configuration
- `WithPort(int)` - Set port
- `WithHost(string)` - Set host

**Logging:**
- `WithLogger(*slog.Logger)` - Custom slog logger
- `WithHTTPLogger(*httplog.Logger)` - HTTP request logger
- `WithLogLevel(slog.Level)` - Set log level

**Middleware:**
- `WithMiddlewareStack(*MiddlewareStack)` - Custom stack
- `WithCORS(*cors.Options)` - Configure CORS
- `WithDefaultCORS()` - CORS with defaults
- `WithHSTS(*gosts.Info)` - Configure HSTS
- `WithDefaultHSTS()` - HSTS with defaults

**Metrics:**
- `WithMetrics(bool)` - Enable metrics (separate mode, backward compatible)
- `WithMetricsCombined()` - Metrics on main app server
- `WithMetricsSeparatePort(port)` - Metrics on separate port
- `WithMetricsPath(path)` - Custom metrics endpoint path
- `WithMetricsMode(mode)` - Set mode explicitly

**Router:**
- `WithRouter(*chi.Mux)` - Use custom router
- `WithHttpin(bool)` - Enable httpin integration

## Project Structure

```
.
├── app/              - Reusable application framework
│   ├── app.go        - Core App initialization
│   ├── middleware.go - Middleware stack system
│   ├── server.go     - Server lifecycle
│   ├── options.go    - Functional options
│   ├── config.go     - Configuration types
│   ├── logging.go    - Logger factories
│   ├── routes.go     - Route helpers
│   └── version.go    - Version middleware
├── cmd/              - Example applications
└── migrations/       - Database migrations
```

## Generate sqlc 

    sqlc generate

## Database

### Start Postgres database
    
    # Postgres
    docker run -it --name postgres -e POSTGRES_PASSWORD=pwd -p 5432:5432 postgres:14-alpine
    
    # MS SQL SEVER
    docker run -it  --name mssql -e ACCEPT_EULA='Y' -e MSSQL_SA_PASSWORD='mssql1Pwd' -p 1433:1433 mcr.microsoft.com/azure-sql-edge

### Connect to database using psql

    psql -h localhost -p 5432 -U postgres

### Create postgres database

     CREATE Role demo_project WITH PASSWORD 'pwd';
     grant demo_project to postgres;
     CREATE DATABASE demo_project_db ENCODING 'UTF8' OWNER demo_project;
     GRANT ALL PRIVILEGES ON DATABASE demo_project_db TO demo_project;
     ALTER ROLE demo_project WITH LOGIN;

     GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO demo_project;

    # https://stackoverflow.com/questions/22135792/permission-denied-to-create-extension-uuid-ossp
    GRANT CREATE ON DATABASE demo_project_db to demo_project;

### Install golang migrate (https://github.com/golang-migrate/migrate)

    brew install golang-migrate 
    
### Create migration

    # If using postgres
    export DATABASE_URL='postgres://demo_project:pwd@localhost:5432/demo_project_db?sslmode=disable'
    migrate create -ext sql -dir migrations/postgres <create_app_user_table>

    # If using MS SQL Server
    export DATABASE_URL="sqlserver://sa:mssql1Pwd@localhost:1433?databaseName=master&integratedSecurity=false&encrypt=false&trustServerCertificate=true"
    migrate create -ext sql -dir migrations/postgres <create_app_user_table>

### Run migration

    migrate -database ${DATABASE_URL} -path migrations/postgres up
    
### Migration Down

    migrate -database ${DATABASE_URL} -path migrations/postgres down

### Migration information table

    schema_migrations

## Fix Dirty DB Migration

    migrate -database ${DATABASE_URL**  -path migrations/postgres force <migration_file_name>


## Test Params

### Query Params

    curl -i localhost:4000/query\?q=DemoQ
    
### Form Params

    curl -i -F 'name=DemoName' localhost:4000/post
    
### Body Params

    curl  -i "localhost:4000/json" -H 'Content-Type: application/json' -d '{"login":"my_login","password":"my_password"}'
    

### Body Json List Params

    curl  -i "localhost:4000/json/list" -H 'Content-Type: application/json' -d '{"emails":["email1@test.com", "email2@test.com", "email3@example.com"]}'

---

## Recent Updates

### v1.6.0 - Middleware Stack System & Flexible Metrics

**Major Refactoring (2025)** - The app package has been completely refactored for flexibility and best practices:

#### ✨ New Features

**1. Middleware Stack System**
- Full control over middleware ordering
- Add, remove, replace, or reorder any middleware
- Named middleware for easy management
- Type-safe builder pattern
- Example: `DefaultMiddlewareStack().InsertAfter("request-id", "auth", authMW).Build()`

**2. Dual Metrics Mode**
- **Combined Mode**: Metrics on same port as app (simple, one port)
- **Separate Mode**: Metrics on dedicated port (production, security)
- Both modes support custom paths
- Backward compatible with old API

**3. Code Reduction**
- `app.go`: 383 lines → 106 lines (72% reduction)
- Removed all code duplication
- Deleted deprecated middleware package
- Clean separation of concerns

#### 📦 New Files

- `middleware.go` - Complete middleware stack system
- `server.go` - Server lifecycle with graceful shutdown
- `options.go` - All functional options centralized
- `routes.go` - Route registration helpers
- `logging.go` - Logger factory functions

#### 🔄 Migration

All existing code continues to work! The refactoring is 100% backward compatible.

**Old code (still works):**
```go
app := app.DefaultApp()
app.Run()
```

**New capabilities:**
```go
// Custom middleware stack
stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", myAuth).
    Build()

// Flexible metrics
app := app.NewApp(
    app.WithMiddlewareStack(stack),
    app.WithMetricsCombined(),  // NEW: Metrics on same port
)
```