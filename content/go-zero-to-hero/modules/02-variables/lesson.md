# Variables and Types

Go is **statically typed** — every variable has a type known at compile time.

## Declaring variables

Two common ways:

```go
var name string = "Go"
age := 25  // short declaration (inside functions)
```

`:=` infers the type from the value on the right.

## Basic types

| Type | Example |
|------|---------|
| `int` | `42` |
| `float64` | `3.14` |
| `string` | `"hello"` |
| `bool` | `true` |

## Zero values

Uninitialized variables get a **zero value**:

- `int` → `0`
- `string` → `""`
- `bool` → `false`

## Using variables

```go
message := "Hello"
count := 3
fmt.Println(message, count)
```

`Println` accepts multiple values separated by spaces.

---

**Next:** Practice declaring and printing variables.
