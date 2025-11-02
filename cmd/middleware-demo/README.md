# Middleware Stack Customization Example

Demonstrates the powerful middleware stack system in chi-demo. Learn how to customize, reorder, add, and remove middleware.

## Features

- ✅ Add custom middleware at any position
- ✅ Remove unwanted middleware
- ✅ Reorder middleware execution
- ✅ Enable/disable middleware dynamically
- ✅ Replace built-in middleware with custom versions

## Running

```bash
go run main.go
```

The example shows multiple middleware stack configurations.

## Key Concepts

### Default Middleware Stack

The framework provides a sensible default:
```
1. RequestID
2. RealIP
3. Recoverer
4. Version
5. HTTPLogger (disabled by default)
6. CORS (disabled by default)
7. NoCache
8. HSTS (disabled by default)
9. Metrics (disabled by default)
```

### Customization Methods

```go
stack := app.DefaultMiddlewareStack().
    Add("name", middleware).           // Add at end
    Prepend("name", middleware).       // Add at start
    Append("name", middleware).        // Add at end
    InsertBefore("ref", "name", mw).   // Insert before reference
    InsertAfter("ref", "name", mw).    // Insert after reference
    Replace("name", newMiddleware).    // Replace existing
    Remove("name").                    // Remove middleware
    Enable("name").                    // Enable disabled middleware
    Disable("name").                   // Disable middleware
    Build()                            // Build final stack
```

## Common Use Cases

### Add Authentication After Request ID
```go
stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", myAuthMiddleware).
    Build()
```

### Remove HSTS for Local Development
```go
stack := app.DefaultMiddlewareStack().
    Remove("hsts").
    Build()
```

### Add Custom Logger
```go
stack := app.DefaultMiddlewareStack().
    Replace("http-logger", myCustomLogger).
    Build()
```

### Add Rate Limiting
```go
stack := app.DefaultMiddlewareStack().
    InsertBefore("recoverer", "rate-limit", rateLimiter).
    Build()
```

## Middleware Execution Order

Middleware executes in the order defined in the stack:
```
Request → Middleware 1 → Middleware 2 → ... → Handler
        ← Middleware 1 ← Middleware 2 ← ... ← Handler
```

Position matters! For example:
- **Authentication** should come early (after RequestID/RealIP)
- **Rate limiting** should come before expensive operations
- **Recovery** should come early to catch panics from all middleware
- **Logging** typically comes early to log all requests

## Tips

- Use descriptive names for your middleware
- Keep the default stack as a starting point
- Test middleware order thoroughly
- Document why you changed the default order
- Consider using `Disable()` instead of `Remove()` for temporary changes
