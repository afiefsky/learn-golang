package grader_test

import (
	"strings"
	"testing"

	"learn-golang/internal/content"
	"learn-golang/internal/grader"
)

func TestChiHealthHandlerGrader(t *testing.T) {
	code := `package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(` + "`" + `{"status":"ok"}` + "`" + `))
}

func main() {
	r := chi.NewRouter()
	r.Get("/health", healthHandler)
	_ = r
}`

	checks := []content.Check{
		{Type: "go_test", Test: `package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHealth(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/health", healthHandler)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "ok") {
		t.Fatalf("body %q", rec.Body.String())
	}
}`},
	}

	g := grader.NewComposite()
	result := g.Grade(grader.GradeInput{
		Code: code,
		Files: map[string]string{
			"go.mod": "module exercise\n\ngo 1.21\n\nrequire github.com/go-chi/chi/v5 v5.0.12\n",
		},
		Checks: checks,
	})
	if !result.Passed {
		t.Fatalf("expected pass, got %+v", result)
	}
}

func TestQuizNormalizeUnderscore(t *testing.T) {
	got := content.NormalizeAnswer("_")
	if got != "_" {
		t.Fatalf("got %q", got)
	}
}

func TestQuizNormalizeInterface(t *testing.T) {
	got := content.NormalizeAnswer("interface{}")
	if got != "interface{}" {
		t.Fatalf("got %q", got)
	}
}

func TestTextGraderContains(t *testing.T) {
	g := grader.NewTextGrader()
	result := g.Grade("fmt.Println(\"hi\")", []content.Check{{Type: "contains", Value: "fmt.Println"}})
	if !result.Passed {
		t.Fatal(result)
	}
}

func TestCompositeNoChecks(t *testing.T) {
	g := grader.NewComposite()
	result := g.Grade(grader.GradeInput{Code: "x", Checks: nil})
	if result.Passed {
		t.Fatal("expected fail with no checks")
	}
	if !strings.Contains(result.Message, "No checks") {
		t.Fatal(result.Message)
	}
}
