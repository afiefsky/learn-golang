# Concurrency Intro

Go runs concurrent tasks with **goroutines** — lightweight threads managed by the runtime.

## Goroutine

```go
go processPayment(id)
```

## Channels

Communicate between goroutines:

```go
ch := make(chan string)
go func() { ch <- "done" }()
msg := <-ch
```

## WaitGroup (sync)

Wait for goroutines to finish:

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()
```

## Fintech use case

Async notifications, batch settlement workers, rate limiters — but **money logic** still needs transactions and careful design.

---

**Next:** Basic goroutine and channel exercises.
