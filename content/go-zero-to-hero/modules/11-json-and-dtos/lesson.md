# JSON & DTOs

**DTO** (Data Transfer Object) in Go is a **struct** shaped for API input/output — not a special keyword.

## Why separate DTOs in fintech

Never return database models directly:

- Hide internal IDs, password hashes, audit fields
- Stable API contract even if DB schema changes
- Validate input before touching money logic

## JSON tags

```go
type TransferRequest struct {
    ToAccount string `json:"to_account"`
    Amount    int64  `json:"amount"`
    Currency  string `json:"currency,omitempty"`
}
```

## Marshal / Unmarshal

```go
data, err := json.Marshal(req)
err = json.Unmarshal(body, &req)
```

## Pattern

| Layer | Struct |
|-------|--------|
| HTTP request | `CreateAccountRequest` |
| HTTP response | `AccountResponse` |
| Database / domain | `Account` (internal) |

Map domain → DTO in a handler or mapper function.

---

**Next:** Build request/response DTOs and parse JSON.
