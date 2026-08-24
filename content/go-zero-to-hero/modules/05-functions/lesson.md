# Functions

Functions group logic you reuse — handlers, validators, fee calculators.

## Basic function

```go
func add(a int, b int) int {
    return a + b
}
```

Same types can be shortened: `func add(a, b int) int`

## Multiple return values

Very common in Go (especially `value, error`):

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

## Named returns

```go
func fees(amount int) (fee int, total int) {
    fee = amount / 100
    total = amount + fee
    return
}
```

## Fintech use case

`CalculateTransferFee(amount int64) int64` — pure function, easy to unit test.

---

**Next:** Write functions with parameters and returns.
