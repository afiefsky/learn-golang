# Slices & Maps

## Slices (dynamic lists)

```go
accounts := []string{"ACC001", "ACC002"}
accounts = append(accounts, "ACC003")
fmt.Println(len(accounts))
```

## Maps (key → value)

```go
balances := map[string]int64{
    "ACC001": 1000000,
    "ACC002": 500000,
}
balances["ACC001"] = 900000
amount, ok := balances["ACC999"]
if !ok {
    fmt.Println("account not found")
}
```

Always check `ok` when reading maps — missing keys return zero value without error.

## Fintech use case

- **Slice:** list of transactions in a statement
- **Map:** account number → balance cache, currency code → rate

---

**Next:** Build and query slices and maps.
