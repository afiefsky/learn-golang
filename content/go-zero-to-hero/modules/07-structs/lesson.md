# Structs & Methods

Structs group related fields — the building block of **DTOs** and domain models.

## Define a struct

```go
type Account struct {
    Number  string
    Balance int64
    Active  bool
}
```

## Create and access

```go
acc := Account{Number: "ACC001", Balance: 1000000, Active: true}
fmt.Println(acc.Number)
```

## Methods (receiver)

```go
func (a Account) IsActive() bool {
    return a.Active
}
```

Pointer receiver when mutating:

```go
func (a *Account) Debit(amount int64) {
    a.Balance -= amount
}
```

## Fintech use case

`Account`, `Transfer`, `Customer` — later you'll separate **domain struct** from **API DTO**.

---

**Next:** Define structs and methods.
