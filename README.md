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
    server := app.DefaultApp()  // CORS, HSTS, logging enabled; metrics disabled
    server.Run()
}
```

### Custom Configuration

```go
myApp := app.NewApp(
    app.WithPort(8080),
    app.WithDefaultCORS(),
    app.WithMetrics(true),  // Enable metrics (combined mode by default)
)

myApp.R.Get("/hello", handleHello)
myApp.Run()
```

### Metrics Modes

**Combined Mode** (default, simple, one port):
```go
// Metrics at http://localhost:3000/metrics (same server)
app := app.NewApp(
    app.WithPort(3000),
    app.WithMetrics(true),  // Combined mode by default
)
```

**Separate Mode** (production, security):
```go
// Metrics at http://localhost:9090/metrics (separate server)
app := app.NewApp(
    app.WithPort(3000),
    app.WithMetricsSeparate(),  // Separate server with default port
)

// Or with custom port
app := app.NewApp(
    app.WithPort(3000),
    app.WithMetricsSeparatePort(9091),
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
- `WithMetrics(bool)` - Enable metrics (combined mode by default)
- `WithMetricsSeparate()` - Enable metrics on separate server (default port 9090)
- `WithMetricsSeparatePort(port)` - Enable metrics on separate server with custom port
- `WithMetricsPath(path)` - Custom metrics endpoint path (works for both modes)
- `WithMetricsMode(mode)` - Set mode explicitly ("combined" or "separate")

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
├── migrations/       - Database migrations
└── dbconn/          - Database connection utilities
```

## Environment Variables

```bash
# Application
APP_ENV=dev              # "dev" or "production"
HOST=localhost
PORT=3000
USE_HTTPIN=false

# Metrics (disabled by default, combined mode when enabled)
METRICS_ENABLED=false    # Set to true to enable
METRICS_MODE=combined    # "combined" or "separate" (default: combined)
METRICS_PATH=/metrics    # Endpoint path
METRICS_HOST=localhost   # Only used in separate mode
METRICS_PORT=9090        # Only used in separate mode
```

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
- **Combined Mode**: Metrics on same port as app (default, simple, one port)
- **Separate Mode**: Metrics on dedicated port (production, security)
- Both modes support custom paths
- Metrics disabled by default (opt-in)

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

// Flexible metrics (combined mode by default)
app := app.NewApp(
    app.WithMiddlewareStack(stack),
    app.WithMetrics(true),  // Metrics on same port by default
)

// Or use separate server for production
app := app.NewApp(
    app.WithMiddlewareStack(stack),
    app.WithMetricsSeparate(),  // Metrics on dedicated port
)
```