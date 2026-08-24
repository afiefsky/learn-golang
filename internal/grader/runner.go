package grader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"learn-golang/internal/content"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Grade(input GradeInput) Result {
	for _, check := range input.Checks {
		if check.Type != "go_test" {
			continue
		}

		result := r.runGoTest(input, check)
		if !result.Passed {
			return result
		}
	}

	return Result{Passed: true, Message: "All checks passed!"}
}

func (r *Runner) runGoTest(input GradeInput, check content.Check) Result {
	dir, err := os.MkdirTemp("", "learn-golang-*")
	if err != nil {
		return Result{Passed: false, Message: "Could not create temp directory.", Errors: []string{err.Error()}}
	}
	defer os.RemoveAll(dir)

	entry := input.EntryFile
	if entry == "" {
		entry = "main.go"
	}

	for name, body := range input.Files {
		if name == entry {
			continue
		}
		if err := writeFile(dir, name, body); err != nil {
			return Result{Passed: false, Message: "Could not write support file.", Errors: []string{err.Error()}}
		}
	}

	if err := writeFile(dir, entry, input.Code); err != nil {
		return Result{Passed: false, Message: "Could not write code.", Errors: []string{err.Error()}}
	}

	testFile := check.TestFile
	if testFile == "" {
		testFile = "main_test.go"
	}
	if strings.TrimSpace(check.Test) != "" {
		if err := writeFile(dir, testFile, check.Test); err != nil {
			return Result{Passed: false, Message: "Could not write tests.", Errors: []string{err.Error()}}
		}
	}

	if _, err := os.Stat(filepath.Join(dir, "go.mod")); os.IsNotExist(err) {
		modContent := "module exercise\n\ngo 1.21\n"
		if err := writeFile(dir, "go.mod", modContent); err != nil {
			return Result{Passed: false, Message: "Could not write go.mod.", Errors: []string{err.Error()}}
		}
	}

	if err := r.tidyModule(dir); err != nil {
		return Result{Passed: false, Message: "Could not resolve modules.", Errors: []string{err.Error()}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "-timeout", "10s", "./...")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(output))

	if err != nil {
		msg := out
		if msg == "" {
			msg = err.Error()
		}
		return Result{
			Passed:  false,
			Message: "Tests failed — check the errors below.",
			Errors:  splitOutput(msg),
		}
	}

	return Result{Passed: true, Message: "Tests passed!"}
}

func (r *Runner) tidyModule(dir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func writeFile(dir, name, body string) error {
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

func splitOutput(output string) []string {
	lines := strings.Split(output, "\n")
	var errors []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			errors = append(errors, line)
		}
	}
	if len(errors) == 0 {
		return []string{output}
	}
	if len(errors) > 20 {
		errors = append(errors[:20], fmt.Sprintf("... and %d more lines", len(errors)-20))
	}
	return errors
}
