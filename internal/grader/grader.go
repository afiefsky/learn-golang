package grader

import (
	"learn-golang/internal/content"
)

type GradeInput struct {
	Code      string
	EntryFile string
	Files     map[string]string
	Checks    []content.Check
}

type Result struct {
	Passed  bool
	Message string
	Errors  []string
}

type Grader interface {
	Grade(input GradeInput) Result
}

type Composite struct {
	text   *TextGrader
	runner *Runner
}

func NewComposite() *Composite {
	return &Composite{
		text:   NewTextGrader(),
		runner: NewRunner(),
	}
}

func (c *Composite) Grade(input GradeInput) Result {
	var textChecks, goTestChecks []content.Check
	for _, check := range input.Checks {
		if check.Type == "go_test" {
			goTestChecks = append(goTestChecks, check)
		} else {
			textChecks = append(textChecks, check)
		}
	}

	if len(textChecks) > 0 {
		result := c.text.Grade(input.Code, textChecks)
		if !result.Passed {
			return result
		}
	}

	if len(goTestChecks) > 0 {
		return c.runner.Grade(input)
	}

	if len(textChecks) > 0 {
		return Result{Passed: true, Message: "All checks passed!"}
	}

	return Result{Passed: false, Message: "No checks configured for this exercise."}
}
