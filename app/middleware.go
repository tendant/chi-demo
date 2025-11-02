package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Middleware represents a chi-compatible middleware function
type Middleware func(http.Handler) http.Handler

// MiddlewareItem represents a named middleware in the stack
type MiddlewareItem struct {
	Name       string
	Middleware Middleware
	Enabled    bool
}

// MiddlewareStack represents an ordered, configurable collection of middleware
type MiddlewareStack struct {
	items []MiddlewareItem
}

// MiddlewareStackBuilder provides a fluent API for building/modifying middleware stacks
type MiddlewareStackBuilder struct {
	items []MiddlewareItem
}

// NewMiddlewareStack creates a new empty middleware stack builder
func NewMiddlewareStack() *MiddlewareStackBuilder {
	return &MiddlewareStackBuilder{
		items: make([]MiddlewareItem, 0),
	}
}

// DefaultMiddlewareStack returns a builder with the recommended middleware stack.
// Order matters! Middleware are applied in the order they're added.
//
// Default stack:
//   - request-id: Injects a request ID into the context
//   - real-ip: Sets a http.Request's RemoteAddr to either X-Forwarded-For or X-Real-IP
//   - recoverer: Recovers from panics, logs the panic, and returns a HTTP 500 status
//   - version: Adds version information to response headers
//   - http-logger: HTTP request/response logger (enabled via WithHTTPLogger)
//   - cors: Cross-Origin Resource Sharing (enabled via WithCORS)
//   - no-cache: Sets response headers to prevent clients from caching
//   - hsts: HTTP Strict Transport Security headers (enabled via WithHSTS)
//   - metrics: Prometheus metrics collection (enabled via WithMetrics)
func DefaultMiddlewareStack() *MiddlewareStackBuilder {
	return NewMiddlewareStack().
		// Core request tracking
		Add("request-id", middleware.RequestID).
		Add("real-ip", middleware.RealIP).
		Add("recoverer", middleware.Recoverer).

		// Version tracking
		Add("version", Version(Commit)).

		// Logging (placeholder - will be configured via WithHTTPLogger)
		AddIf("http-logger", nil, false).

		// CORS (placeholder - will be configured via WithCORS)
		AddIf("cors", nil, false).

		// Security
		Add("no-cache", middleware.NoCache).

		// HSTS (placeholder - will be configured via WithHSTS)
		AddIf("hsts", nil, false).

		// Metrics (placeholder - will be configured via WithMetrics)
		AddIf("metrics", nil, false)
}

// MinimalMiddlewareStack returns just the essential middleware.
// Useful as a starting point for building custom stacks.
func MinimalMiddlewareStack() *MiddlewareStackBuilder {
	return NewMiddlewareStack().
		Add("request-id", middleware.RequestID).
		Add("recoverer", middleware.Recoverer)
}

// Add appends a middleware to the end of the stack
func (b *MiddlewareStackBuilder) Add(name string, mw Middleware) *MiddlewareStackBuilder {
	b.items = append(b.items, MiddlewareItem{
		Name:       name,
		Middleware: mw,
		Enabled:    true,
	})
	return b
}

// AddIf conditionally adds middleware based on the enabled flag
func (b *MiddlewareStackBuilder) AddIf(name string, mw Middleware, enabled bool) *MiddlewareStackBuilder {
	b.items = append(b.items, MiddlewareItem{
		Name:       name,
		Middleware: mw,
		Enabled:    enabled,
	})
	return b
}

// Prepend adds a middleware to the beginning of the stack
func (b *MiddlewareStackBuilder) Prepend(name string, mw Middleware) *MiddlewareStackBuilder {
	b.items = append([]MiddlewareItem{{
		Name:       name,
		Middleware: mw,
		Enabled:    true,
	}}, b.items...)
	return b
}

// Append is an alias for Add (for clarity in code)
func (b *MiddlewareStackBuilder) Append(name string, mw Middleware) *MiddlewareStackBuilder {
	return b.Add(name, mw)
}

// InsertBefore inserts middleware before the named middleware.
// If the target is not found, appends to the end.
func (b *MiddlewareStackBuilder) InsertBefore(targetName, name string, mw Middleware) *MiddlewareStackBuilder {
	idx := b.findIndex(targetName)
	if idx == -1 {
		// Target not found, append to end
		return b.Add(name, mw)
	}

	item := MiddlewareItem{Name: name, Middleware: mw, Enabled: true}
	b.items = append(b.items[:idx], append([]MiddlewareItem{item}, b.items[idx:]...)...)
	return b
}

// InsertAfter inserts middleware after the named middleware.
// If the target is not found, appends to the end.
func (b *MiddlewareStackBuilder) InsertAfter(targetName, name string, mw Middleware) *MiddlewareStackBuilder {
	idx := b.findIndex(targetName)
	if idx == -1 {
		// Target not found, append to end
		return b.Add(name, mw)
	}

	item := MiddlewareItem{Name: name, Middleware: mw, Enabled: true}
	b.items = append(b.items[:idx+1], append([]MiddlewareItem{item}, b.items[idx+1:]...)...)
	return b
}

// Replace replaces a named middleware with a new one.
// If not found, does nothing.
func (b *MiddlewareStackBuilder) Replace(name string, mw Middleware) *MiddlewareStackBuilder {
	idx := b.findIndex(name)
	if idx != -1 {
		b.items[idx].Middleware = mw
	}
	return b
}

// Remove removes a named middleware from the stack completely.
// If not found, does nothing.
func (b *MiddlewareStackBuilder) Remove(name string) *MiddlewareStackBuilder {
	idx := b.findIndex(name)
	if idx != -1 {
		b.items = append(b.items[:idx], b.items[idx+1:]...)
	}
	return b
}

// Enable enables a named middleware (keeps it in stack).
// Useful for re-enabling disabled middleware.
func (b *MiddlewareStackBuilder) Enable(name string) *MiddlewareStackBuilder {
	idx := b.findIndex(name)
	if idx != -1 {
		b.items[idx].Enabled = true
	}
	return b
}

// Disable disables a named middleware (keeps it in stack but won't apply).
// Useful for temporarily disabling middleware without removing it.
func (b *MiddlewareStackBuilder) Disable(name string) *MiddlewareStackBuilder {
	idx := b.findIndex(name)
	if idx != -1 {
		b.items[idx].Enabled = false
	}
	return b
}

// Clear removes all middleware from the stack
func (b *MiddlewareStackBuilder) Clear() *MiddlewareStackBuilder {
	b.items = make([]MiddlewareItem, 0)
	return b
}

// Build returns the final middleware stack
func (b *MiddlewareStackBuilder) Build() *MiddlewareStack {
	// Make a copy of items to prevent external modification
	items := make([]MiddlewareItem, len(b.items))
	copy(items, b.items)
	return &MiddlewareStack{items: items}
}

// List returns the names of all middleware in order with enabled status.
// Format: "✓ name" for enabled, "✗ name" for disabled
func (b *MiddlewareStackBuilder) List() []string {
	names := make([]string, 0, len(b.items))
	for _, item := range b.items {
		status := "✓"
		if !item.Enabled {
			status = "✗"
		}
		names = append(names, fmt.Sprintf("%s %s", status, item.Name))
	}
	return names
}

// findIndex finds the index of a named middleware
func (b *MiddlewareStackBuilder) findIndex(name string) int {
	for i, item := range b.items {
		if item.Name == name {
			return i
		}
	}
	return -1
}

// Apply applies all enabled middleware to a chi router in order
func (s *MiddlewareStack) Apply(r chi.Router) {
	for _, item := range s.items {
		if item.Enabled && item.Middleware != nil {
			r.Use(item.Middleware)
		}
	}
}

// Items returns a copy of the middleware items
func (s *MiddlewareStack) Items() []MiddlewareItem {
	items := make([]MiddlewareItem, len(s.items))
	copy(items, s.items)
	return items
}

// List returns the names of all middleware in order with enabled status.
// Format: "✓ name" for enabled, "✗ name" for disabled
func (s *MiddlewareStack) List() []string {
	names := make([]string, 0, len(s.items))
	for _, item := range s.items {
		status := "✓"
		if !item.Enabled {
			status = "✗"
		}
		names = append(names, fmt.Sprintf("%s %s", status, item.Name))
	}
	return names
}

// Clone creates a builder from an existing stack
func (s *MiddlewareStack) Clone() *MiddlewareStackBuilder {
	items := make([]MiddlewareItem, len(s.items))
	copy(items, s.items)
	return &MiddlewareStackBuilder{items: items}
}
