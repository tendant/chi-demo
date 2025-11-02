# Examples

Practical, copy-paste-ready examples for common use cases. Each example is designed to be a starting point for your project.

## Quick Start

### 🚀 **[basic/](basic/)** - Start Here!
The simplest setup. Copy this to start your project.
```go
app := app.DefaultApp()
app.R.Get("/", handler)
app.Run()
```
**What you get**: CORS, HSTS, logging, graceful shutdown, panic recovery

## Core Features

### 🔧 **[middleware-demo/](middleware-demo/)** - Customize Middleware
Learn to add, remove, reorder, and replace middleware.
```go
stack := app.DefaultMiddlewareStack().
    InsertAfter("request-id", "auth", myAuth).
    Remove("hsts").
    Build()
```
**Use when**: You need custom middleware ordering

### 📥 **[input-handling/](input-handling/)** - Parse Inputs
Handle query parameters, forms, and JSON with type safety.
```go
type Input struct {
    Name string `in:"query=name"`
    Age  int    `in:"query=age"`
}
```
**Use when**: Processing user input from HTTP requests

### 📝 **[gen/](gen/)** - Simple CRUD API
In-memory REST API with Create, Read, Update, Delete operations.
```bash
curl -X POST http://localhost:3000/tasks \
  -d '{"title":"Task 1"}'
```
**Use when**: Building REST APIs with basic CRUD operations

## Authentication

### 🔐 **[apikey/](apikey/)** - API Key Auth
Simple API key authentication with SHA256 hashing.
```go
curl http://localhost:3000/protected \
  -H "X-API-KEY: secret-key"
```
**Use when**: Microservices, mobile backends, CLI tools

### 🛡️ **[auth-examples/](auth-examples/)** - Advanced Auth
Multiple authentication patterns:
- **jwt-auth/** - JWT with social login support
- **oidc-keycloak/** - OpenID Connect integration
- **oidc-custom/** - Custom OIDC provider mapping

**Use when**: User authentication with enterprise SSO or social logins

## Documentation & Tooling

### 📚 **[openapi/](openapi/)** - API Documentation
Serve OpenAPI specs with Swagger UI for interactive API docs.
```
http://localhost:3000/swagger/
```
**Use when**: Building public APIs or need interactive documentation

## Learning & Comparison

### 📖 **[first/](first/)** - Raw Chi
Minimal chi router setup with **no framework**. See what chi-demo abstracts away.
```go
r := chi.NewRouter()
r.Get("/", handler)
http.ListenAndServe(":3000", r)
```
**Use when**: Learning chi basics or need minimal dependencies

## How to Use These Examples

1. **Browse**: Look at the examples above
2. **Read**: Each example has a detailed README
3. **Run**: `cd cmd/example-name && go run main.go`
4. **Copy**: Copy the code to your project
5. **Customize**: Adapt to your needs

## Decision Guide

**I want to...**

- **Start a new project** → [basic/](basic/)
- **Add authentication** → [apikey/](apikey/) or [auth-examples/](auth-examples/)
- **Parse form data** → [input-handling/](input-handling/)
- **Build a REST API** → [gen/](gen/)
- **Customize middleware** → [middleware-demo/](middleware-demo/)
- **Add API docs** → [openapi/](openapi/)
- **Learn raw chi** → [first/](first/)

## Example Complexity

```
Simplicity ←──────────────────────────→ Power

first     basic   apikey   gen   middleware   input     auth      openapi
  ●        ●        ●       ●        ●         ●         ●          ●
  │        │        │       │        │         │         │          │
Raw     Start   Simple   CRUD   Customize   Forms   Enterprise  Docs
Chi     Here    Auth     API    Middleware  Input    Auth       +UI
```

## Common Patterns Not Covered

These examples focus on HTTP server patterns. For other needs:

- **Database**: Use your preferred ORM (gorm, sqlc, etc.)
- **Testing**: Standard Go testing with `httptest`
- **WebSockets**: Use gorilla/websocket with chi
- **Background Jobs**: Use your preferred worker library
- **gRPC**: Chi is for HTTP; use grpc-go separately
- **Rate Limiting**: Use chi middleware like tollbooth
- **Caching**: Use your preferred caching library

## Contributing Examples

Have a useful pattern? Consider contributing:
1. Focus on a single concept
2. Keep it simple and well-documented
3. Make it copy-paste ready
4. Include a README with testing commands
5. Follow the existing example structure

## Next Steps

After reviewing examples:
1. Check the main [README.md](../README.md) for framework documentation
2. Review the [app/](../app/) package for available options
3. Start building your application!
