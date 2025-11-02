# Input Handling Example

Demonstrates how to parse and validate different types of HTTP input using the `ggicci/httpin` library integrated with the chi-demo app framework.

## Features

- **Query parameters**: Simple string/int parsing
- **Array query parameters**: Multiple values for same key
- **Boolean query parameters**: Auto-conversion
- **Form data**: application/x-www-form-urlencoded
- **JSON body**: application/json parsing

## Running

```bash
go run main.go
```

## Testing

### Query Parameters
```bash
# Simple query
curl "http://localhost:3000/query?name=John&age=30"

# Array and boolean query
curl "http://localhost:3000/users?is_member=true&age_range=18&age_range=60"
```

### Form Data
```bash
curl -X POST http://localhost:3000/form \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=John&description=Test user"
```

### JSON Body
```bash
curl -X POST http://localhost:3000/json \
  -H "Content-Type: application/json" \
  -d '{"name":"John","description":"Test user","tags":["admin","developer"]}'
```

## Key Concepts

1. **httpin integration**: Enable with `app.WithHttpin(true)`
2. **Middleware pattern**: Use `httpin.NewInput()` as middleware
3. **Context extraction**: Get parsed input from `r.Context().Value(httpin.Input)`
4. **Type safety**: Define structs with `in:` tags for validation

## Learn More

- [httpin documentation](https://ggicci.github.io/httpin/)
- [Directive reference](https://ggicci.github.io/httpin/directives/query)
