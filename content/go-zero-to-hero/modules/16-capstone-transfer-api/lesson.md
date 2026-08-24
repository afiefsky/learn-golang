# Capstone — Mini Transfer API

Build this **outside** the in-browser editor as a separate folder (e.g. `capstone/transfer-api/`) — this is your **portfolio piece** for Indonesian fintech interviews.

## Requirements

```
POST /v1/accounts
GET  /v1/accounts/{id}
POST /v1/transfers
```

Headers:

- `Authorization: Bearer <token>` on protected routes
- `Idempotency-Key: <uuid>` on transfers

## Suggested stack

- **Chi** router
- **DTOs** for request/response
- **SQLite** (local) or **PostgreSQL** (closer to production)
- **JWT** or simple API key for learning
- Structured logging (`log/slog`)

## Architecture

```
HTTP Handler → Service → Repository → DB
```

## Demo script (for interviews)

1. Create account
2. Login / get token
3. Transfer with idempotency key
4. Replay same key — no double debit
5. Show error envelope on insufficient balance

## Reference projects

- Tech School Simple Bank
- GoBank (gRPC + gateway — advanced)

Check off each item in the **project checklist** below when done.

---

There are no coding exercises in this module — complete the checklist, then take the short quiz.
