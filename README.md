# Learn Go — Local Codecademy-Style App

Personal, local-first Go learning app with interactive exercises, quizzes, and progress tracking — **Fintech Edition** tuned for Indonesian banking/fintech engineering paths.

## Requirements

- Go 1.21+ (with `go` on your PATH for exercise grading with `go test`)

## Run

```bash
go run ./cmd/server
```

Open [http://localhost:8080](http://localhost:8080)

Optional: set `ADDR=:3000` to change the port.

## Learning path (16 modules)

**Fundamentals (1–10):** Go syntax → structs → errors → interfaces → concurrency

**Industry bridge (11–16):**

| Module | Topic |
|--------|--------|
| 11 | JSON & DTOs |
| 12 | Chi REST API |
| 13 | REST design for fintech |
| 14 | PostgreSQL & transactions |
| 15 | Auth, JWT & middleware |
| 16 | Capstone transfer API (project checklist) |

## Structure

```
cmd/server/          HTTP server (Chi)
internal/api/        REST handlers
internal/content/    YAML/Markdown/JSON loader
internal/grader/     Text checks + go test runner (multi-file)
internal/progress/   SQLite progress store
content/             Course curriculum
web/                 Bootstrap frontend
data/                SQLite DB (created at runtime)
```

## Adding content

Edit files under `content/go-zero-to-hero/modules/<module-id>/`:

- `module.yaml` — title, description, optional `project` checklist
- `lesson.md` — lesson narrative (Markdown)
- `exercises.yaml` — coding exercises and grader checks
- `quiz.json` — quiz questions

Register new modules in `content/go-zero-to-hero/course.yaml`.

### Exercise checks

```yaml
checks:
  - type: contains
    value: 'fmt.Println'
  - type: regex
    pattern: 'func main\(\)'
  - type: go_test
    test: |
      package main
      import "testing"
      func TestX(t *testing.T) { ... }
files:
  go.mod: |
    module exercise
    go 1.21
    require github.com/go-chi/chi/v5 v5.0.12
```

## Capstone

Module 16 guides building a **Mini Transfer API** portfolio project (Chi + DTOs + DB + auth) in a separate folder — ideal for fintech interviews.

## License

Personal learning project.
