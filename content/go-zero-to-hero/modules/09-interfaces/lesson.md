# Interfaces

An **interface** is a set of method signatures. Types satisfy interfaces **implicitly** — no `implements` keyword.

## Example

```go
type BalanceReader interface {
    GetBalance(accountNumber string) (int64, error)
}
```

Any type with matching methods implements the interface.

## Why fintech teams use them

- Swap real DB for mock in tests
- Multiple implementations (cache vs postgres)
- Handler depends on interface, not concrete struct

```go
type AccountService struct {
    repo BalanceReader
}
```

## Small interfaces

Go culture prefers **small, focused interfaces** (often 1–2 methods).

---

**Next:** Define and implement an interface.
