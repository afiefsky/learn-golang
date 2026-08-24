# PostgreSQL & Transactions

Production Go fintech backends commonly use **PostgreSQL** + **sqlc** or careful SQL.

## Connect

```go
import "database/sql"
_ "github.com/jackc/pgx/v5/stdlib"

db, err := sql.Open("pgx", connString)
```

## Query with context

```go
row := db.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id=$1", id)
err := row.Scan(&balance)
```

## Transaction (transfer)

```go
tx, err := db.BeginTx(ctx, nil)
// debit source
// credit destination
return tx.Commit()
// on any error: tx.Rollback()
```

## Rules for money

1. **Always** use transactions for debit+credit
2. Check `rowsAffected` / balance in same txn
3. Use `idempotency_key` unique constraint

## Local Postgres (optional)

```bash
docker run -d --name pg -e POSTGRES_PASSWORD=secret -p 5432:5432 postgres:16
```

## sqlc (industry tool)

Write SQL → generate type-safe Go — used in Simple Bank / GoBank tutorials.

---

**Next:** Transaction pattern exercises (in-memory simulation).
