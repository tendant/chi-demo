# Authentication Examples

This directory contains multiple authentication patterns for different use cases. Pick the one that matches your requirements and copy it to your project.

## Examples

### 1. JWT with go-pkgz/auth (`jwt-auth/`)
Full-featured authentication service with:
- Direct login provider (username/password)
- JWT token generation and validation
- XSRF protection
- Cookie-based sessions
- Support for OAuth providers (GitHub, Google, etc.)

**Use when**: You need a complete auth solution with social logins

### 2. OIDC with Keycloak (`oidc-keycloak/`)
OpenID Connect integration:
- OAuth2 authorization code flow
- State and nonce validation
- Token exchange and verification
- Works with Keycloak, Auth0, Okta, etc.

**Use when**: Integrating with enterprise identity providers

## See Also

- **cmd/apikey** - Simple API key authentication for REST APIs
- **cmd/secure** - Alternative OIDC setup with custom provider mapping

## Quick Comparison

| Feature | JWT (go-pkgz/auth) | OIDC | API Key |
|---------|-------------------|------|---------|
| Social Login | ✅ Yes | ✅ Yes | ❌ No |
| Enterprise SSO | ⚠️ Limited | ✅ Yes | ❌ No |
| Simplicity | ⚠️ Medium | ⚠️ Medium | ✅ Simple |
| Token Type | JWT | JWT/OIDC | SHA256 Hash |
| Setup Complexity | Medium | Medium-High | Low |
| External Dependencies | None | Identity Provider | None |

## Notes

- All examples can be integrated with the chi-demo app framework
- See individual README files in each subdirectory for detailed setup instructions
- Remember to change secrets and keys before deploying to production
