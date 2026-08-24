# Control Flow

Programs need to make decisions. Go provides `if`, `else`, and `switch`.

## if / else

```go
score := 85
if score >= 60 {
    fmt.Println("Pass")
} else {
    fmt.Println("Fail")
}
```

Go requires `{` on the same line as `if`. No parentheses around the condition.

## if with short statement

```go
if n := 10; n > 5 {
    fmt.Println("n is big")
}
```

## switch

```go
day := "Mon"
switch day {
case "Mon":
    fmt.Println("Start of week")
case "Fri":
    fmt.Println("Almost weekend")
default:
    fmt.Println("Midweek")
}
```

Unlike many languages, Go `switch` does **not** fall through by default.

---

**Next:** Write programs that branch based on conditions.
