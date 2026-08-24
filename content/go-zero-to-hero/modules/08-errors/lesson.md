# Error Handling

Go has no exceptions. Failures are **`error` values** returned alongside results.

## The idiom

```go
result, err := doSomething()
if err != nil {
    return err
}
```

## Creating errors

```go
import "errors"
import "fmt"

errors.New("insufficient balance")
fmt.Errorf("transfer failed: account %s", id)
```

## Custom error types (optional)

```go
type ValidationError struct {
    Field string
    Msg   string
}
func (e ValidationError) Error() string { return e.Msg }
```

## Fintech rule

**Never ignore errors** in money paths. Log and return appropriate HTTP status later.

---

**Next:** Return and check errors.
