# Auth, JWT & Middleware

Fintech APIs protect money endpoints with **authentication** and **authorization**.

## JWT (JSON Web Token) — common pattern

Three parts: header.payload.signature

Claims example: `sub` (user id), `exp` (expiry)

Flow:

1. `POST /v1/login` → access token + refresh token
2. Client sends `Authorization: Bearer <token>`
3. Middleware validates and attaches user to `context`

## Chi middleware

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
        if token == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        // validate token, set ctx
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## RBAC

- **Customer** — own accounts only
- **Admin** — support operations (audited)

## Security rules

- Never log tokens or passwords
- Short-lived access tokens + refresh rotation
- Some banks use PASETO instead of JWT — same middleware idea

---

**Next:** Build auth middleware patterns.
