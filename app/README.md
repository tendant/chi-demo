# Chi Demo App Package

A flexible, type-safe Go web application framework built on top of [chi router](https://github.com/go-chi/chi), featuring a powerful middleware stack system with full customization capabilities.

## Features

- ✅ **Functional Options Pattern** - Type-safe configuration
- ✅ **Flexible Middleware Stack** - Full control over middleware ordering
- ✅ **Default Sensible Defaults** - Works out of the box
- ✅ **Graceful Shutdown** - Proper signal handling
- ✅ **Structured Logging** - slog integration with environment-based configuration
- ✅ **CORS Support** - Configurable cross-origin resource sharing
- ✅ **Prometheus Metrics** - Optional metrics collection
- ✅ **Security Headers** - HSTS support
- ✅ **Configuration via Environment** - 12-factor app compliant

## Quick Start

### Simplest Example

```go
package main

import "github.com/tendant/chi-demo/app"

func main() {
    myApp := app.NewApp()
    myApp.Run()
}
```

### With Custom Port

```go
myApp := app.NewApp(
    app.WithPort(8080),
)
myApp.Run()
```

### With Common Features

```go
myApp := app.NewApp(
    app.WithPort(8080),
    app.WithDefaultCORS(),
    app.WithMetrics(true),
)

// Add your routes
myApp.R.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World!"))
})

myApp.Run()
```

### Use Recommended Defaults

```go
// Includes: logging, CORS, metrics, HSTS, httpin
myApp := app.DefaultApp()

// Add your routes
app.RegisterHealthzRoutes(myApp.R)
myApp.R.Get("/api/users", handleUsers)

myApp.Run()
```

## Middleware Stack

The middleware stack system provides full control over middleware ordering and configuration while maintaining type safety.

### Default Middleware Stack

The default stack includes (in order):

1. `request-id` - Injects unique request ID
2. `real-ip` - Extracts real IP from X-Forwarded-For
3. `recoverer` - Panic recovery
4. `version` - Version info in headers
5. `http-logger` - HTTP request/response logging (optional)
6. `cors` - Cross-origin resource sharing (optional)
7. `no-cache` - Cache control headers
8. `hsts` - HTTP Strict Transport Security (optional)
9. `metrics` - Prometheus metrics (optional)

### Customizing Middleware Stack

```go
// Start with default stack and customize
stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", myAuthMiddleware).
    Remove("hsts").
    Replace("http-logger", myCustomLogger).
    Disable("metrics").
    Build()

myApp := app.NewApp(
    app.WithMiddlewareStack(stack),
)
```

### Building from Scratch

```go
stack := app.NewMiddlewareStack().
    Add("request-id", middleware.RequestID).
    Add("my-auth", authMiddleware).
    Add("logger", loggerMiddleware).
    Build()

myApp := app.NewApp(
    app.WithMiddlewareStack(stack),
)
```

### Middleware Stack Methods

#### Builder Methods

- `Add(name, middleware)` - Append to end
- `Prepend(name, middleware)` - Add to beginning
- `Append(name, middleware)` - Alias for Add
- `InsertBefore(target, name, middleware)` - Insert before named middleware
- `InsertAfter(target, name, middleware)` - Insert after named middleware
- `Replace(name, middleware)` - Replace existing middleware
- `Remove(name)` - Remove middleware completely
- `Enable(name)` - Enable disabled middleware
- `Disable(name)` - Disable without removing
- `Clear()` - Remove all middleware
- `Build()` - Finalize and return stack

#### Inspection Methods

- `List()` - Get middleware names with status (✓/✗)

### Minimal Stack

For lightweight applications:

```go
stack := app.MinimalMiddlewareStack() // Just request-id + recoverer
    .Append("logger", myLogger).
    Build()
```

## Configuration

### Environment Variables

```bash
# Application
APP_ENV=dev              # "dev" or "production"
HOST=localhost
PORT=3000
USE_HTTPIN=false

# Metrics
METRICS_ENABLED=true
METRICS_HOST=localhost
METRICS_PORT=9090
```

### Functional Options

#### Configuration

- `WithConfig(AppConfig)` - Set complete configuration
- `WithPort(int)` - Set port
- `WithHost(string)` - Set host

#### Logging

- `WithLogger(*slog.Logger)` - Custom slog logger
- `WithHTTPLogger(*httplog.Logger)` - HTTP request logger
- `WithLogLevel(slog.Level)` - Set log level

#### Middleware

- `WithMiddlewareStack(*MiddlewareStack)` - Custom stack
- `WithCORS(*cors.Options)` - Configure CORS
- `WithDefaultCORS()` - CORS with defaults
- `WithMetrics(bool)` - Enable Prometheus metrics
- `WithHSTS(*gosts.Info)` - Configure HSTS
- `WithDefaultHSTS()` - HSTS with defaults

#### Router

- `WithRouter(*chi.Mux)` - Use custom router
- `WithHttpin(bool)` - Enable httpin integration

## Examples

### Custom Authentication Middleware

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !validateToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", authMiddleware).
    Build()

myApp := app.NewApp(
    app.WithMiddlewareStack(stack),
    app.WithDefaultCORS(),
)
```

### Production Configuration

```go
stack := app.DefaultMiddlewareStack().
    InsertAfter("recoverer", "auth", authMiddleware).
    InsertAfter("auth", "rbac", rbacMiddleware).
    Build()

myApp := app.NewApp(
    app.WithConfig(app.DefaultAppConfig()),
    app.WithLogger(productionLogger),
    app.WithMiddlewareStack(stack),
    app.WithCORS(strictCORSOptions),
    app.WithMetrics(true),
    app.WithDefaultHSTS(),
)

// Register routes
app.RegisterHealthzRoutes(myApp.R)
app.RegisterVersionRoutes(myApp.R)

myApp.R.Route("/api", func(r chi.Router) {
    r.Get("/users", handleUsers)
    r.Post("/users", createUser)
})

myApp.Run()
```

### Conditional Middleware

```go
stack := app.DefaultMiddlewareStack()

if os.Getenv("ENABLE_AUTH") == "true" {
    stack.InsertAfter("request-id", "auth", authMiddleware)
}

if os.Getenv("ENABLE_RATE_LIMIT") == "true" {
    stack.Append("rate-limit", rateLimitMiddleware)
}

myApp := app.NewApp(
    app.WithMiddlewareStack(stack.Build()),
)
```

### Debug Middleware Stack

```go
stack := app.DefaultMiddlewareStack()

fmt.Println("Middleware stack:")
for _, name := range stack.List() {
    fmt.Println(" ", name)
}

myApp := app.NewApp(
    app.WithMiddlewareStack(stack.Build()),
)
```

## Architecture

### File Structure

```
app/
├── app.go          - Core App struct, NewApp(), Run()
├── config.go       - Configuration types and validation
├── server.go       - Server lifecycle management
├── middleware.go   - Middleware stack system
├── logging.go      - Logger factories
├── options.go      - Functional options
├── routes.go       - Default route helpers
└── version.go      - Version middleware
```

### Key Types

```go
type App struct {
    R               *chi.Mux
    Config          AppConfig
    Logger          *slog.Logger
    HTTPLogger      *httplog.Logger
    middlewareStack *MiddlewareStack
    // ...
}

type MiddlewareStack struct {
    items []MiddlewareItem
}

type MiddlewareItem struct {
    Name       string
    Middleware Middleware
    Enabled    bool
}
```

## Migration Guide

### From Old API

**Before:**
```go
app := app.Default()
app.Run()
```

**After:**
```go
app := app.DefaultApp()
app.Run()
```

**Before (custom config):**
```go
app := app.NewApp(
    WithAppConfig(config),
    WithMetrics(true),
    WithCors(corsOpts),
)
```

**After:**
```go
app := app.NewApp(
    app.WithConfig(config),
    app.WithMetrics(true),
    app.WithCORS(corsOpts),
)
```

## Best Practices

1. **Use `DefaultApp()` for most cases** - It includes sensible defaults
2. **Customize via options** - Use functional options for configuration
3. **Order matters** - Middleware execute in order they're added
4. **Auth early** - Add authentication middleware early in the stack
5. **Logging after auth** - Log after authentication for security
6. **Metrics at the start** - Capture all requests including failures
7. **Use environment variables** - Follow 12-factor app principles

## Testing

```go
func TestApp(t *testing.T) {
    myApp := app.NewApp(
        app.WithPort(0), // Random port for testing
    )

    // Add test routes
    myApp.R.Get("/test", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    // Test with httptest
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    myApp.R.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

## License

Same as parent project.
