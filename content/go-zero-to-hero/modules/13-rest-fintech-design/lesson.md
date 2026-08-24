# REST Design for Fintech

API design matters as much as code in Indonesian banking/fintech teams.

## Versioning

```
POST /v1/transfers
GET  /v1/accounts/{id}
```

Never break mobile apps — add `/v2` instead of changing `/v1`.

## Status codes (common)

| Code | When |
|------|------|
| 200 | OK, GET success |
| 201 | Resource created |
| 400 | Bad JSON / validation |
| 401 | Missing/invalid auth |
| 409 | Duplicate / conflict |
| 422 | Business rule failed |
| 500 | Unexpected server error |

## Idempotency-Key

Clients retry on network failure. Same key = same result, no double transfer:

```
POST /v1/transfers
Idempotency-Key: 7b2c-4f91-a003
```

Store key + response; return cached response on replay.

## Error envelope

```json
{
  "code": "INSUFFICIENT_BALANCE",
  "message": "Saldo tidak mencukupi",
  "details": []
}
```

## Logging (audit)

Log: request ID, user ID, endpoint, outcome. **Never** log passwords, OTP, full card numbers.

---

**Next:** Choose correct responses for fintech scenarios.
