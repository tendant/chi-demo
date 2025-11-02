# Raw Chi Router Example

The absolute minimum chi router setup with **no chi-demo framework**. This example shows what raw chi looks like for comparison.

## Purpose

This example helps you understand:
- What raw chi router looks like
- What the chi-demo framework abstracts away
- When you might want to use raw chi vs the framework

## Running

```bash
go run main.go
```

## Testing

```bash
curl http://localhost:3000/
```

## Comparison

### Raw Chi (this example)
```go
r := chi.NewRouter()
r.Get("/", handler)
http.ListenAndServe(":3000", r)
```

**What you don't get:**
- No CORS
- No security headers
- No structured logging
- No graceful shutdown
- No panic recovery
- No middleware stack management

### With chi-demo Framework
```go
app := app.DefaultApp()
app.R.Get("/", handler)
app.Run()
```

**What you get automatically:**
- ✅ CORS with defaults
- ✅ HSTS security headers
- ✅ Structured logging (slog)
- ✅ Graceful shutdown
- ✅ Panic recovery
- ✅ Request ID tracking
- ✅ Flexible middleware stack
- ✅ Environment-based config

## When to Use Raw Chi

Use raw chi when:
- Learning chi basics
- Building a microservice with minimal dependencies
- You need absolute control over every aspect
- The framework abstractions don't fit your use case

## When to Use chi-demo Framework

Use the framework when:
- Building production APIs
- You want sensible defaults
- You need common patterns (auth, logging, metrics)
- You want to focus on business logic, not boilerplate
- You're building multiple services and want consistency
