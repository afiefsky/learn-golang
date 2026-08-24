# Chi REST API

This learning app itself uses **Chi** — see [`internal/api/handlers.go`](../../../../internal/api/handlers.go) in the repo.

## Router basics

```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Get("/health", healthHandler)
r.Post("/v1/transfers", createTransfer)
http.ListenAndServe(":8080", r)
```

## URL parameters

```go
r.Get("/accounts/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
})
```

## Handler signature

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ok"}`))
}
```

## Middleware

```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
```

## Layering (fintech style)

`Handler` → decode DTO → call `Service` → map to response DTO → write JSON

---

**Next:** Write Chi handlers and routes.
