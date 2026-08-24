# Getting Started with Go

Welcome to Go! In this module you'll write and run your first program.

## What is Go?

Go (often called **Golang**) is an open-source language created at Google. It is known for:

- **Simplicity** — small language, easy to read
- **Speed** — compiled to native machine code
- **Built-in concurrency** — goroutines and channels (later modules)

## Your first program

Every Go program starts with a **package** declaration. Executable programs use `package main`.

The entry point is always:

```go
func main() {
    // your code here
}
```

## Printing output

Use the `fmt` package (short for "format"):

```go
import "fmt"

fmt.Println("Hello, World!")
```

`Println` prints a line and adds a newline automatically.

## Running Go code

Save a file as `main.go`, then run:

```bash
go run main.go
```

Later you'll learn `go build` to compile a binary.

---

**Next:** Complete the exercises below to practice writing your first Go programs.
