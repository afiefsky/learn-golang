package grader

import (
	"fmt"
	"regexp"
	"strings"

	"learn-golang/internal/content"
)

type TextGrader struct{}

func NewTextGrader() *TextGrader {
	return &TextGrader{}
}

func (g *TextGrader) Grade(code string, checks []content.Check) Result {
	var errors []string

	for _, check := range checks {
		msg := check.Message
		if msg == "" {
			msg = defaultMessage(check)
		}

		switch check.Type {
		case "contains":
			if !strings.Contains(code, check.Value) {
				errors = append(errors, msg)
			}
		case "not_contains":
			if strings.Contains(code, check.Value) {
				errors = append(errors, msg)
			}
		case "regex":
			re, err := regexp.Compile(check.Pattern)
			if err != nil {
				errors = append(errors, fmt.Sprintf("invalid regex check: %v", err))
				continue
			}
			if !re.MatchString(code) {
				errors = append(errors, msg)
			}
		case "equals":
			if strings.TrimSpace(code) != strings.TrimSpace(check.Value) {
				errors = append(errors, msg)
			}
		default:
			errors = append(errors, fmt.Sprintf("unknown check type: %s", check.Type))
		}
	}

	if len(errors) > 0 {
		return Result{
			Passed:  false,
			Message: "Not quite — keep trying!",
			Errors:  errors,
		}
	}

	return Result{Passed: true, Message: "All checks passed!"}
}

func defaultMessage(check content.Check) string {
	switch check.Type {
	case "contains":
		return fmt.Sprintf("Your code should include: %s", check.Value)
	case "not_contains":
		return fmt.Sprintf("Your code should not include: %s", check.Value)
	case "regex":
		return fmt.Sprintf("Your code should match pattern: %s", check.Pattern)
	case "equals":
		return "Your code does not match the expected solution yet."
	default:
		return "Check failed."
	}
}
