# Basic Example

The simplest way to get started with chi-demo. Perfect for beginners!

## What You Get

- ✅ **CORS** enabled with permissive defaults
- ✅ **HSTS** security headers
- ✅ **Structured logging** with slog (pretty-printed in dev)
- ✅ **Graceful shutdown** on SIGINT/SIGTERM
- ✅ **Request ID** tracking
- ✅ **Panic recovery**

## Running

```bash
go run main.go
```

## Testing

```bash
# Root endpoint
curl http://localhost:3000/

# Dynamic route with path parameter
curl http://localhost:3000/hello/World
```

## Copy & Paste to Start Your Project

This example is designed to be copied directly. Just:

1. Copy this `main.go` to your project
2. Add your routes after line 14
3. Run it!

## Customization

Want to customize? See other examples:
- **cmd/middleware-demo** - Customize the middleware stack
- **cmd/apikey** - Add API key authentication
- **cmd/input-handling** - Parse query/form/JSON inputs

## Next Steps

```go
// Enable metrics
myApp := app.NewApp(
    app.WithDefaultCORS(),
    app.WithMetrics(true),  // Metrics at http://localhost:3000/metrics
)

// Change port
myApp := app.NewApp(
    app.WithPort(8080),
)

// Customize middleware
stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", myAuthMiddleware).
    Build()

myApp := app.NewApp(
    app.WithMiddlewareStack(stack),
)
```
