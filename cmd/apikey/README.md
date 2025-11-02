# API Key Authentication Example

Simple API key authentication for REST APIs using SHA256-hashed keys.

## Features

- ✅ API key validation via `X-API-KEY` header
- ✅ SHA256 hashing for secure key storage
- ✅ Protected and public routes
- ✅ Integration with chi-demo app framework

## Running

```bash
go run main.go
```

## Testing

```bash
# Public endpoint (no auth required)
curl http://localhost:3000/

# Protected endpoint (requires API key)
curl http://localhost:3000/protected \
  -H "X-API-KEY: secret-key-1"

# Will return 401 without valid key
curl http://localhost:3000/protected
```

## How It Works

1. API keys are stored as SHA256 hashes
2. Incoming keys are hashed and compared
3. Valid keys allow access to protected routes
4. Invalid/missing keys return 401 Unauthorized

## Adding Your Own Keys

```go
apikeys := app.ApiKeyConfig{
    ApiKeys: map[string][]byte{
        // Generate hash: echo -n "your-secret-key" | sha256sum
        "user1": []byte("hash-of-key-1"),
        "user2": []byte("hash-of-key-2"),
    },
    HeaderKey: "X-API-KEY", // or use "Authorization"
}
```

## Production Considerations

- Store API key hashes in a database, not in code
- Use strong, random API keys (32+ characters)
- Consider rate limiting per API key
- Log API key usage for auditing
- Rotate keys periodically
- Use HTTPS in production

## Use Cases

- ✅ Microservice-to-microservice authentication
- ✅ Third-party API integrations
- ✅ Mobile app backends
- ✅ CLI tool authentication
- ❌ Not for user-facing web apps (use JWT/OIDC instead)
