# Loops & Range

Go has only one loop keyword: **`for`**. There is no `while` — you use `for` with a condition instead.

## Classic for loop

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

## While-style loop

```go
n := 10
for n > 0 {
    n--
}
```

## Range over slice

```go
amounts := []int{100, 250, 75}
for i, amount := range amounts {
    fmt.Println(i, amount)
}
```

Use `_` to ignore index or value:

```go
for _, amount := range amounts {
    fmt.Println(amount)
}
```

## Fintech use case

Loop through transaction rows, statement lines, or paginated API results.

---

**Next:** Practice summing values and iterating slices.
